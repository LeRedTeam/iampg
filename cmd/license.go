package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LeRedTeam/iampg/license"
	"github.com/spf13/cobra"
)

var licenseCmd = &cobra.Command{
	Use:    "license",
	Short:  "License management commands",
	Hidden: true, // Hidden from normal users
}

var licenseStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current license status",
	RunE:  runLicenseStatus,
}

var licenseGenerateKeypairCmd = &cobra.Command{
	Use:   "generate-keypair",
	Short: "Generate a new keypair for license signing (admin only)",
	RunE:  runGenerateKeypair,
}

var licenseGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a license key (admin only)",
	RunE:  runGenerateLicense,
}

var (
	licenseEmail      string
	licenseTier       string
	licenseValidDays  int
	licensePrivateKey string
)

func init() {
	rootCmd.AddCommand(licenseCmd)
	licenseCmd.AddCommand(licenseStatusCmd)
	licenseCmd.AddCommand(licenseGenerateKeypairCmd)
	licenseCmd.AddCommand(licenseGenerateCmd)

	licenseGenerateCmd.Flags().StringVar(&licenseEmail, "email", "", "License holder email")
	licenseGenerateCmd.Flags().StringVar(&licenseTier, "tier", "pro", "License tier (pro, team)")
	licenseGenerateCmd.Flags().IntVar(&licenseValidDays, "days", 365, "Validity period in days")
	licenseGenerateCmd.Flags().StringVar(&licensePrivateKey, "private-key", "", "Private key for signing (or IAMPG_PRIVATE_KEY env)")

	licenseGenerateCmd.MarkFlagRequired("email")
}

func runLicenseStatus(cmd *cobra.Command, args []string) error {
	lic, err := license.Current()
	if err != nil {
		return fmt.Errorf("license error: %w", err)
	}

	fmt.Printf("Tier:    %s\n", lic.Tier)
	if lic.IsPaid() {
		fmt.Printf("Email:   %s\n", lic.Email)
		fmt.Printf("Expires: %s\n", lic.ExpiresAt.Format("2006-01-02"))
	}

	fmt.Println("\nFeatures:")
	features := []string{"run", "parse", "json", "yaml", "terraform", "refine", "sarif", "enforce", "diff"}
	for _, f := range features {
		status := "✗"
		if lic.HasFeature(f) {
			status = "✓"
		}
		fmt.Printf("  %s %s\n", status, f)
	}

	return nil
}

func runGenerateKeypair(cmd *cobra.Command, args []string) error {
	keypair, err := license.GenerateKeyPair()
	if err != nil {
		return err
	}

	output, _ := json.MarshalIndent(keypair, "", "  ")
	fmt.Println(string(output))

	fmt.Fprintln(os.Stderr, "\n⚠️  Store the private_key securely. It cannot be recovered.")
	fmt.Fprintln(os.Stderr, "   Embed the public_key in the binary at build time.")

	return nil
}

func runGenerateLicense(cmd *cobra.Command, args []string) error {
	privateKey := licensePrivateKey
	if privateKey == "" {
		privateKey = os.Getenv("IAMPG_PRIVATE_KEY")
	}
	if privateKey == "" {
		return fmt.Errorf("private key required: use --private-key or IAMPG_PRIVATE_KEY")
	}

	tier := license.Tier(licenseTier)
	if tier != license.TierPro && tier != license.TierTeam {
		return fmt.Errorf("invalid tier: must be 'pro' or 'team'")
	}

	key, err := license.GenerateLicenseKey(privateKey, licenseEmail, tier, licenseValidDays)
	if err != nil {
		return err
	}

	fmt.Println(key)

	return nil
}
