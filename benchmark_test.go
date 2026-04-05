package smtpmock

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkMessagesCopy(b *testing.B) {
	for _, count := range []int{0, 1, 10, 100, 1000} {
		messages := &messages{}
		for i := 0; i < count; i++ {
			messages.append(createNotEmptyMessage())
		}

		b.Run(messageCountName(count), func(b *testing.B) {
			for b.Loop() {
				messages.copy()
			}
		})
	}
}

func BenchmarkMessagesCopyConcurrent(b *testing.B) {
	for _, count := range []int{1, 10, 100} {
		messages := &messages{}
		for i := 0; i < count; i++ {
			messages.append(createNotEmptyMessage())
		}

		b.Run(messageCountName(count), func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					messages.copy()
				}
			})
		})
	}
}

func BenchmarkMessagesPurge(b *testing.B) {
	for _, count := range []int{1, 10, 100} {
		b.Run(messageCountName(count), func(b *testing.B) {
			for b.Loop() {
				b.StopTimer()
				messages := &messages{}
				for i := 0; i < count; i++ {
					messages.append(createNotEmptyMessage())
				}
				b.StartTimer()

				messages.purge()
			}
		})
	}
}

func BenchmarkMessagesAppend(b *testing.B) {
	messages := &messages{}
	message := createNotEmptyMessage()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			messages.append(message)
		}
	})
}

func BenchmarkSMTPSessionMinimal(b *testing.B) {
	_, hostAddress, portNumber := newBenchServer(b)

	b.ResetTimer()
	for b.Loop() {
		runMinimalSMTPSession(hostAddress, portNumber) //nolint:errcheck
	}
}

func BenchmarkSMTPSessionFull(b *testing.B) {
	_, hostAddress, portNumber := newBenchServer(b)

	b.ResetTimer()
	for b.Loop() {
		runFullSMTPBenchSession(hostAddress, portNumber) //nolint:errcheck
	}
}

func BenchmarkSMTPSessionConcurrent(b *testing.B) {
	_, hostAddress, portNumber := newBenchServer(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runFullSMTPBenchSession(hostAddress, portNumber) //nolint:errcheck
		}
	})
}

func BenchmarkServerMessages(b *testing.B) {
	for _, count := range []int{1, 10, 100} {
		server := New(ConfigurationAttr{})
		for i := 0; i < count; i++ {
			server.messages.append(createNotEmptyMessage())
		}

		b.Run(messageCountName(count), func(b *testing.B) {
			for b.Loop() {
				server.Messages()
			}
		})
	}
}

func BenchmarkServerMessagesAndPurge(b *testing.B) {
	server := New(ConfigurationAttr{})

	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		for i := 0; i < 10; i++ {
			server.messages.append(createNotEmptyMessage())
		}
		b.StartTimer()

		server.MessagesAndPurge()
	}
}

func BenchmarkServerMessagesConcurrentReadWrite(b *testing.B) {
	server := New(ConfigurationAttr{})
	for i := 0; i < 10; i++ {
		server.messages.append(createNotEmptyMessage())
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			server.Messages()
			server.messages.append(createNotEmptyMessage())
		}
	})
}

func BenchmarkWaitForMessages(b *testing.B) {
	server := New(ConfigurationAttr{})
	for i := 0; i < 5; i++ {
		server.messages.append(createNotEmptyMessage())
	}

	b.ResetTimer()
	for b.Loop() {
		server.WaitForMessages(5, time.Second) //nolint:errcheck
	}
}

func BenchmarkWaitForMessagesWithArrival(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		server := New(ConfigurationAttr{})
		b.StartTimer()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(5 * time.Millisecond)
			for i := 0; i < 3; i++ {
				server.messages.append(createNotEmptyMessage())
			}
		}()

		server.WaitForMessages(3, time.Second) //nolint:errcheck
		wg.Wait()
	}
}

func BenchmarkRegexMatchCommand(b *testing.B) {
	server := New(ConfigurationAttr{})
	commands := []string{"HELO example.com", "MAIL FROM:<user@test.com>", "RCPT TO:<user@test.com>", "DATA", "RSET", "NOOP", "QUIT"}

	b.ResetTimer()
	for b.Loop() {
		for _, command := range commands {
			server.isInvalidCmd(command)
		}
	}
}

func BenchmarkRegexValidateHelo(b *testing.B) {
	request := "EHLO mail.example.com"
	b.ResetTimer()
	for b.Loop() {
		validHeloComplexCmdRegex.MatchString(request)
	}
}

func BenchmarkRegexValidateMailfrom(b *testing.B) {
	request := "MAIL FROM:<sender@example.com>"
	b.ResetTimer()
	for b.Loop() {
		validMailfromComplexCmdRegex.MatchString(request)
		regexCaptureGroupCompiled(request, validMailfromComplexCmdRegex, 2)
	}
}
