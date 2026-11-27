package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/tekkamanendless/constantcontact"
)

func main() {
	ctx := context.Background()

	constantContactClient := constantcontact.NewClient()

	oauth2Input := constantcontact.PostOauth2TokenInput{
		ClientID:     os.Getenv("CONSTANT_CONTACT_CLIENT_ID"),
		RefreshToken: os.Getenv("CONSTANT_CONTACT_REFRESH_TOKEN"),
		GrantType:    constantcontact.GrantTypeRefreshToken,
		Scope:        "contact_data campaign_data offline_access",
	}
	err := constantContactClient.AuthenticateWithOauth2(ctx, oauth2Input)
	if err != nil {
		panic(fmt.Sprintf("Could not refresh Constant Contact OAuth2 token: %v", err))
	}

	/*
		{
			_, err := constantContactClient.ActivitiesContactDelete(ctx, constantcontact.ActivitiesContactDeleteInput{
				ContactIDs: []string{"bfec67cc-e7d5-11f1-ba1b-02420a320002"},
			})
			if err != nil {
				panic(fmt.Sprintf("Could not delete Constant Contact contacts: %v", err))
			}

			os.Exit(1)
		}
	*/
	/*
		{
			contact, err := constantContactClient.GetContact(ctx, constantcontact.GetContactInput{
				ContactID: "3fb7be7c-12be-11f1-9f3d-0242ac2f7082",
			})
			if err != nil {
				panic(fmt.Sprintf("Could not get Constant Contact contact: %v", err))
			}
			fmt.Printf("Constant Contact contact: %+v\n", contact)

			contact.EmailAddress.PermissionToSend = "unsubscribed"
			contact.EmailAddress.OptOutReason = "Test of API on 2026-11-24."
			contact.UpdateSource = "Contact"

			_, err = constantContactClient.PutContact(ctx, constantcontact.PutContactInput(contact))
			if err != nil {
				panic(fmt.Sprintf("Could not put Constant Contact contact: %v", err))
			}

			os.Exit(1)
		}
	*/

	stdinContents, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(fmt.Sprintf("Could not read stdin: %v", err))
	}

	contactIDStatusMap := map[string]bool{}
	contactIDEmailAddressMap := map[string]string{}
	contactIDOptOutReasonMap := map[string]string{}
	contactIDOptOutDateMap := map[string]string{}

	if len(stdinContents) > 0 {
		csvReader := csv.NewReader(strings.NewReader(string(stdinContents)))
		rows, err := csvReader.ReadAll()
		if err != nil {
			panic(fmt.Sprintf("Could not read CSV: %v", err))
		}
		if len(rows) == 0 {
			panic("No rows provided")
		}

		header := rows[0]
		rows = rows[1:]

		headerMap := map[string]int{}
		for i, column := range header {
			headerMap[column] = i
		}

		if _, ok := headerMap["email_address"]; !ok {
			panic("No email address column provided")
		}
		if _, ok := headerMap["contact_id"]; !ok {
			panic("No contact ID column provided")
		}

		for _, row := range rows {
			contactID := row[headerMap["contact_id"]]
			emailAddress := row[headerMap["email_address"]]
			contactIDStatusMap[contactID] = false
			contactIDEmailAddressMap[contactID] = emailAddress
			contactIDOptOutReasonMap[contactID] = row[headerMap["opt_out_reason"]]
			contactIDOptOutDateMap[contactID] = row[headerMap["opt_out_date"]]
		}
	} else {
		getEmailsOutput, err := constantContactClient.GetEmails(ctx, constantcontact.GetEmailsInput{
			Limit: 500,
		})
		if err != nil {
			panic(fmt.Sprintf("Could not get Constant Contact emails: %v", err))
		}
		slices.SortFunc(getEmailsOutput.Campaigns, func(a, b constantcontact.GetEmailsOutputCampaign) int {
			return strings.Compare(a.CreatedAt, b.CreatedAt)
		})

		fmt.Printf("Constant Contact emails: %d\n", len(getEmailsOutput.Campaigns))

		for _, email := range getEmailsOutput.Campaigns {
			fmt.Printf("Email: %s\n", email.Name)

			emailOutput, err := constantContactClient.GetEmailID(ctx, constantcontact.GetEmailIDInput{
				CampaignID: email.CampaignID,
			})
			if err != nil {
				panic(fmt.Sprintf("Could not get Constant Contact email: %v", err))
			}

			fmt.Printf("Email activities: %d\n", len(emailOutput.CampaignActivities))

			for _, activity := range emailOutput.CampaignActivities {
				fmt.Printf("Activity: %s\n", activity.CampaignActivityID)

				{
					emailSendsReportOutput, err := constantContactClient.GetEmailSendsReport(ctx, constantcontact.GetEmailSendsReportInput{
						CampaignActivityID: activity.CampaignActivityID,
					})
					if err != nil {
						panic(fmt.Sprintf("Could not get Constant Contact email sends report: %v", err))
					}

					fmt.Printf("Email sends: %d\n", len(emailSendsReportOutput.TrackingActivities))

					for _, trackingActivity := range emailSendsReportOutput.TrackingActivities {
						contactIDStatusMap[trackingActivity.ContactID] = true
						contactIDEmailAddressMap[trackingActivity.ContactID] = trackingActivity.EmailAddress
						delete(contactIDOptOutReasonMap, trackingActivity.ContactID)
						delete(contactIDOptOutDateMap, trackingActivity.ContactID)
					}
				}

				{
					emailOptoutsReportOutput, err := constantContactClient.GetEmailOptoutsReport(ctx, constantcontact.GetEmailOptoutsReportInput{
						CampaignActivityID: activity.CampaignActivityID,
					})
					if err != nil {
						panic(fmt.Sprintf("Could not get Constant Contact email sends report: %v", err))
					}

					fmt.Printf("Email unsubscribes: %d\n", len(emailOptoutsReportOutput.TrackingActivities))

					for _, trackingActivity := range emailOptoutsReportOutput.TrackingActivities {
						contactIDStatusMap[trackingActivity.ContactID] = false
						contactIDEmailAddressMap[trackingActivity.ContactID] = trackingActivity.EmailAddress
						contactIDOptOutReasonMap[trackingActivity.ContactID] = trackingActivity.OptOutReason
						contactIDOptOutDateMap[trackingActivity.ContactID] = trackingActivity.CreatedTime
					}
				}
			}
		}
	}

	contactIDsToUnsubscribe := []string{}
	for contactID, status := range contactIDStatusMap {
		fmt.Printf("Email address: %s\n", contactIDEmailAddressMap[contactID])
		if status {
			continue
		}

		if reason, ok := contactIDOptOutReasonMap[contactID]; ok {
			fmt.Printf("Opt out reason: %s\n", reason)
		}
		contactIDsToUnsubscribe = append(contactIDsToUnsubscribe, contactID)
	}

	fmt.Printf("Contact IDs to unsubscribe: %d\n", len(contactIDsToUnsubscribe))

	rows := [][]string{}
	rows = append(rows, []string{"contact_id", "email_address", "opt_out_reason", "opt_out_date"})
	for _, contactID := range contactIDsToUnsubscribe {
		rows = append(rows, []string{contactID, contactIDEmailAddressMap[contactID], contactIDOptOutReasonMap[contactID], contactIDOptOutDateMap[contactID]})
	}

	fmt.Printf("\n")
	csvWriter := csv.NewWriter(os.Stdout)
	csvWriter.WriteAll(rows)
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		panic(fmt.Sprintf("Could not write CSV: %v", err))
	}
	fmt.Printf("\n")

	for _, contactID := range contactIDsToUnsubscribe {
		fmt.Printf("Contact ID: %s\n", contactID)
		fmt.Printf("Email address: %s\n", contactIDEmailAddressMap[contactID])
		fmt.Printf("Opt out reason: %s\n", contactIDOptOutReasonMap[contactID])
		fmt.Printf("Opt out date: %s\n", contactIDOptOutDateMap[contactID])

		contact, err := constantContactClient.GetContact(ctx, constantcontact.GetContactInput{
			ContactID: contactID,
		})
		if err != nil {
			panic(fmt.Sprintf("Could not get Constant Contact contact: %v", err))
		}
		fmt.Printf("Constant Contact contact: %+v\n", contact)

		if contact.EmailAddress.PermissionToSend == "unsubscribed" {
			fmt.Printf("Contact already unsubscribed: %s\n", contact.EmailAddress.OptOutReason)
			continue
		}

		contact.EmailAddress.PermissionToSend = "unsubscribed"
		contact.EmailAddress.OptOutReason = contactIDOptOutReasonMap[contactID]
		contact.UpdateSource = "Contact"

		_, err = constantContactClient.PutContact(ctx, constantcontact.PutContactInput(contact))
		if err != nil {
			panic(fmt.Sprintf("Could not put Constant Contact contact: %v", err))
		}
	}
}
