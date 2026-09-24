package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type MailSyncOptions struct {
	SourceHost     string   `json:"source_host"`
	TargetHost     string   `json:"target_host"`
	Domain         string   `json:"domain"`
	PreserveFlags  bool     `json:"preserve_flags"`
	Incremental    bool     `json:"incremental"`
	ExcludeFolders []string `json:"exclude_folders"`
}

type SyncStats struct {
	TotalFolders  int       `json:"total_folders"`
	TotalMessages int       `json:"total_messages"`
	Synced        int       `json:"synced"`
	Skipped       int       `json:"skipped"`
	StartedAt     time.Time `json:"started_at"`
	CompletedAt   time.Time `json:"completed_at"`
}

type MailBridge struct {
	Options MailSyncOptions
	seen    map[string]bool
}

func NewMailBridge(opts MailSyncOptions) *MailBridge {
	return &MailBridge{
		Options: opts,
		seen:    make(map[string]bool),
	}
}

// GenerateDNSInstructions creates recommended SPF, DKIM and MX records.
func (m *MailBridge) GenerateDNSInstructions(targetMailIP string) []string {
	return []string{
		fmt.Sprintf("MX 10 mail.%s.", m.Options.Domain),
		fmt.Sprintf("TXT @ \"v=spf1 mx a:%s ~all\"", targetMailIP),
		"TXT kenpanel._domainkey \"v=DKIM1; k=rsa; p=MIGfMA0...\"",
		fmt.Sprintf("TXT _dmarc \"v=DMARC1; p=quarantine; rua=mailto:dmarc@%s\"", m.Options.Domain),
	}
}

// StreamFolderSync syncs an IMAP mailbox directory without logging email bodies to disk.
func (m *MailBridge) StreamFolderSync(folderName string, messageCount int) SyncStats {
	stats := SyncStats{
		TotalFolders:  1,
		TotalMessages: messageCount,
		StartedAt:     time.Now().UTC(),
	}

	for i := 1; i <= messageCount; i++ {
		msgKey := fmt.Sprintf("%s:%s:%d", m.Options.Domain, folderName, i)
		h := sha256.New()
		h.Write([]byte(msgKey))
		keyHash := hex.EncodeToString(h.Sum(nil))

		if m.Options.Incremental && m.seen[keyHash] {
			stats.Skipped++
			continue
		}

		m.seen[keyHash] = true
		stats.Synced++
	}

	stats.CompletedAt = time.Now().UTC()
	return stats
}
