package smtpmock

import (
	"fmt"
	"net"
	"net/smtp"
	"sync"
	"testing"
	"time"
)

// --- Message copy benchmarks (directly relevant to issue #200) ---

// BenchmarkMessagesCopy benchmarks copying messages with varying slice sizes.
// This is the hot path in fetchMessages() that issue #200 optimizes.
func BenchmarkMessagesCopy(b *testing.B) {
	for _, count := range []int{0, 1, 10, 100, 1000} {
		msgs := &messages{}
		for i := 0; i < count; i++ {
			msgs.append(createNotEmptyMessage())
		}

		b.Run(messageCountName(count), func(b *testing.B) {
			for b.Loop() {
				msgs.copy()
			}
		})
	}
}

// BenchmarkMessagesCopyConcurrent benchmarks concurrent reads on the messages slice.
func BenchmarkMessagesCopyConcurrent(b *testing.B) {
	for _, count := range []int{1, 10, 100} {
		msgs := &messages{}
		for i := 0; i < count; i++ {
			msgs.append(createNotEmptyMessage())
		}

		b.Run(messageCountName(count), func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					msgs.copy()
				}
			})
		})
	}
}

// BenchmarkMessagesPurge benchmarks purge (copy + clear) with varying sizes.
func BenchmarkMessagesPurge(b *testing.B) {
	for _, count := range []int{1, 10, 100} {
		b.Run(messageCountName(count), func(b *testing.B) {
			for b.Loop() {
				b.StopTimer()
				msgs := &messages{}
				for i := 0; i < count; i++ {
					msgs.append(createNotEmptyMessage())
				}
				b.StartTimer()

				msgs.purge()
			}
		})
	}
}

// BenchmarkMessagesAppend benchmarks appending messages concurrently.
func BenchmarkMessagesAppend(b *testing.B) {
	msgs := &messages{}
	msg := createNotEmptyMessage()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			msgs.append(msg)
		}
	})
}

// --- SMTP session benchmarks ---

// runMinimalSMTPSession runs a HELO + QUIT session, returns error on failure.
func runMinimalSMTPSession(host string, port int) error {
	conn, err := net.DialTimeout(networkProtocol, serverWithPortNumber(host, port), 2*time.Second)
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	if err := client.Hello("bench.test"); err != nil {
		return err
	}
	return client.Quit()
}

// runFullSMTPBenchSession runs a full SMTP session for benchmarking.
func runFullSMTPBenchSession(host string, port int) error {
	conn, err := net.DialTimeout(networkProtocol, serverWithPortNumber(host, port), 2*time.Second)
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	if err := client.Hello("bench.test"); err != nil {
		return err
	}

	sender, receiver := "sender@bench.test", "receiver@bench.test"
	if err := client.Mail(sender); err != nil {
		return err
	}
	if err := client.Rcpt(receiver); err != nil {
		return err
	}
	wc, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := wc.Write(messageBody(sender, receiver)); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return client.Quit()
}

// newBenchServer creates and starts a server for benchmarking.
func newBenchServer(b *testing.B) (*Server, string, int) {
	b.Helper()
	server := New(ConfigurationAttr{})
	if err := server.Start(); err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() { server.Stop() }) //nolint:errcheck
	return server, server.configuration.hostAddress, server.PortNumber()
}

// BenchmarkSMTPSessionMinimal benchmarks a minimal SMTP session (HELO + QUIT only).
func BenchmarkSMTPSessionMinimal(b *testing.B) {
	_, host, port := newBenchServer(b)

	b.ResetTimer()
	for b.Loop() {
		runMinimalSMTPSession(host, port) //nolint:errcheck
	}
}

// BenchmarkSMTPSessionFull benchmarks a full SMTP session
// (HELO, MAIL FROM, RCPT TO, DATA, message body, QUIT).
func BenchmarkSMTPSessionFull(b *testing.B) {
	_, host, port := newBenchServer(b)

	b.ResetTimer()
	for b.Loop() {
		runFullSMTPBenchSession(host, port) //nolint:errcheck
	}
}

// BenchmarkSMTPSessionConcurrent benchmarks concurrent full SMTP sessions.
func BenchmarkSMTPSessionConcurrent(b *testing.B) {
	_, host, port := newBenchServer(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runFullSMTPBenchSession(host, port) //nolint:errcheck
		}
	})
}

// --- Server message retrieval benchmarks ---

// BenchmarkServerMessages benchmarks Messages() with accumulated messages.
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

// BenchmarkServerMessagesAndPurge benchmarks MessagesAndPurge() cycle.
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

// BenchmarkServerMessagesConcurrentReadWrite benchmarks concurrent message
// reading while new messages are being appended (simulates real usage).
func BenchmarkServerMessagesConcurrentReadWrite(b *testing.B) {
	server := New(ConfigurationAttr{})
	// Pre-populate some messages
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

// --- fetchMessages / WaitForMessages benchmarks ---

// BenchmarkWaitForMessages benchmarks WaitForMessages with messages already available.
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

// BenchmarkWaitForMessagesWithArrival benchmarks WaitForMessages where messages
// arrive asynchronously (the realistic polling scenario that issue #200 targets).
func BenchmarkWaitForMessagesWithArrival(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		server := New(ConfigurationAttr{})
		b.StartTimer()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Simulate messages arriving with a small delay
			time.Sleep(5 * time.Millisecond)
			for i := 0; i < 3; i++ {
				server.messages.append(createNotEmptyMessage())
			}
		}()

		server.WaitForMessages(3, time.Second) //nolint:errcheck
		wg.Wait()
	}
}

// --- Regex benchmarks ---

// BenchmarkRegexMatchCommand benchmarks the isInvalidCmd check that runs on every SMTP command.
func BenchmarkRegexMatchCommand(b *testing.B) {
	server := New(ConfigurationAttr{})
	commands := []string{"HELO example.com", "MAIL FROM:<user@test.com>", "RCPT TO:<user@test.com>", "DATA", "RSET", "NOOP", "QUIT"}

	b.ResetTimer()
	for b.Loop() {
		for _, cmd := range commands {
			server.isInvalidCmd(cmd)
		}
	}
}

// BenchmarkRegexValidateHelo benchmarks HELO command argument validation.
func BenchmarkRegexValidateHelo(b *testing.B) {
	request := "EHLO mail.example.com"
	b.ResetTimer()
	for b.Loop() {
		validHeloComplexCmdRegex.MatchString(request)
	}
}

// BenchmarkRegexValidateMailfrom benchmarks MAIL FROM argument validation and capture.
func BenchmarkRegexValidateMailfrom(b *testing.B) {
	request := "MAIL FROM:<sender@example.com>"
	b.ResetTimer()
	for b.Loop() {
		validMailfromComplexCmdRegex.MatchString(request)
		regexCaptureGroupCompiled(request, validMailfromComplexCmdRegex, 2)
	}
}

// --- Helper ---

func messageCountName(n int) string {
	if n == 0 {
		return "empty"
	}
	return fmt.Sprintf("%dmsgs", n)
}
