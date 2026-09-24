package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kentralo/kenpanel-mailbridge/engine"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "dns-plan":
		planCmd := flag.NewFlagSet("dns-plan", flag.ExitOnError)
		domain := planCmd.String("domain", "example.com", "Mail domain")
		targetIP := planCmd.String("ip", "203.0.113.10", "Destination server IP")
		planCmd.Parse(os.Args[2:])

		mb := engine.NewMailBridge(engine.MailSyncOptions{Domain: *domain})
		records := mb.GenerateDNSInstructions(*targetIP)

		fmt.Printf("Recommended DNS Records for Cutover (%s):\n\n", *domain)
		for _, r := range records {
			fmt.Printf("  %s\n", r)
		}

	case "version":
		fmt.Println("kenpanel-mailbridge v2.3.0 (standalone engine)")

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: kenpanel-mailbridge <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  dns-plan    Generate required SPF, DKIM, DMARC and MX records before cutover")
	fmt.Println("  sync        Stream mailboxes between source and target IMAP/Maildir servers")
	fmt.Println("  version     Show version")
}
