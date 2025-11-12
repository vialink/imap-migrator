package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// saveReport saves the migration report to a text file.
func saveReport(report MigrationReport) error {
	// Create reports directory if it doesn't exist
	reportsDir := "reports"
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return fmt.Errorf("error creating reports directory: %w", err)
	}
	
	// Generate filename based on email and timestamp
	timestamp := report.StartTime.Format("20060102_150405")
	safeEmail := strings.ReplaceAll(report.SourceEmail, "@", "_at_")
	safeEmail = strings.ReplaceAll(safeEmail, ".", "_")
	filename := fmt.Sprintf("migration_%s_%s.txt", safeEmail, timestamp)
	filePath := filepath.Join(reportsDir, filename)
	
	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating report file: %w", err)
	}
	defer file.Close()
	
	// Write header
	fmt.Fprintf(file, "═══════════════════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(file, "                    IMAP MIGRATION REPORT\n")
	fmt.Fprintf(file, "═══════════════════════════════════════════════════════════════════════════\n\n")
	
	// General information
	fmt.Fprintf(file, "Source:      %s\n", report.SourceEmail)
	fmt.Fprintf(file, "Destination: %s\n", report.DestinationEmail)
	fmt.Fprintf(file, "Start:       %s\n", report.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "End:         %s\n", report.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "Duration:    %s\n", formatDuration(report.Duration))
	
	if report.Success {
		fmt.Fprintf(file, "Status:      ✓ COMPLETED SUCCESSFULLY\n")
	} else {
		fmt.Fprintf(file, "Status:      ✗ INTERRUPTED (see errors below)\n")
	}
	
	fmt.Fprintf(file, "\n")
	
	// Overall summary
	fmt.Fprintf(file, "───────────────────────────────────────────────────────────────────────────\n")
	fmt.Fprintf(file, "                           SUMMARY\n")
	fmt.Fprintf(file, "───────────────────────────────────────────────────────────────────────────\n\n")
	
	fmt.Fprintf(file, "Total folders processed:         %d\n", report.TotalFolders)
	fmt.Fprintf(file, "Total messages at source:        %d\n", report.TotalSourceMsgs)
	fmt.Fprintf(file, "Total messages copied:           %d\n", report.TotalCopied)
	fmt.Fprintf(file, "Total messages failed:           %d\n", report.TotalFailed)
	fmt.Fprintf(file, "Total messages skipped:          %d\n", report.TotalSkipped)
	
	if report.TotalSourceMsgs > 0 {
		successRate := float64(report.TotalCopied) / float64(report.TotalSourceMsgs) * 100
		fmt.Fprintf(file, "Success rate:                    %.2f%%\n", successRate)
	}
	
	fmt.Fprintf(file, "\n")
	
	// Per-folder details
	fmt.Fprintf(file, "───────────────────────────────────────────────────────────────────────────\n")
	fmt.Fprintf(file, "                      FOLDER DETAILS\n")
	fmt.Fprintf(file, "───────────────────────────────────────────────────────────────────────────\n\n")
	
	// Table header
	fmt.Fprintf(file, "%-50s %8s %8s %8s %8s\n", "FOLDER", "SOURCE", "COPIED", "FAILED", "SKIPPED")
	fmt.Fprintf(file, "%-50s %8s %8s %8s %8s\n", strings.Repeat("-", 50), "--------", "--------", "--------", "--------")
	
	for _, folder := range report.Folders {
		// Truncate folder name if too long
		folderName := folder.Name
		if len(folderName) > 50 {
			folderName = folderName[:47] + "..."
		}
		
		fmt.Fprintf(file, "%-50s %8d %8d %8d %8d\n",
			folderName,
			folder.SourceMessages,
			folder.CopiedMessages,
			folder.FailedMessages,
			folder.SkippedMessages)
	}
	
	fmt.Fprintf(file, "\n")
	
	// Errors (if any)
	if len(report.Errors) > 0 {
		fmt.Fprintf(file, "───────────────────────────────────────────────────────────────────────────\n")
		fmt.Fprintf(file, "                            ERRORS\n")
		fmt.Fprintf(file, "───────────────────────────────────────────────────────────────────────────\n\n")
		
		for i, errMsg := range report.Errors {
			fmt.Fprintf(file, "%d. %s\n", i+1, errMsg)
		}
		
		fmt.Fprintf(file, "\n")
	} else {
		fmt.Fprintf(file, "───────────────────────────────────────────────────────────────────────────\n")
		fmt.Fprintf(file, "                  ✓ NO ERRORS RECORDED\n")
		fmt.Fprintf(file, "───────────────────────────────────────────────────────────────────────────\n\n")
	}
	
	// Footer
	fmt.Fprintf(file, "═══════════════════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(file, "Report generated at: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "═══════════════════════════════════════════════════════════════════════════\n")
	
	return nil
}

// formatDuration formats a duration in a readable way.
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	} else {
		return fmt.Sprintf("%ds", seconds)
	}
}
