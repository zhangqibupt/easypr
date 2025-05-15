package cmd

import (
	"context"
	"testing"

	"github.com/fatih/color"
	"github.com/google/go-github/v57/github"
)

func Test_CreateComments(t *testing.T) {
	//prLink := "https://github.freewheel.tv/biz-ui/superadmin/pull/1689"
	//performReview(prLink)

	comment := &github.PullRequestComment{
		Body: github.String("This is the body of a comment."),
		Path: github.String("func/configuration/config/env.prd.us-east-1.yml"),
		//Line:        github.Int(10),
		CommitID:    github.String("91a7b754b0ebd920873edd7d6c40f2058290a388"),
		SubjectType: github.String("file"),

		//InReplyTo: github.Int(0), // Optional, set to 0 for a new comment
	}

	if err := initClient(context.Background()); err != nil {
		color.Red("Failed to init GitHub client: %s", err)
		return
	}

	_, _, err := client.PullRequests.CreateComment(context.Background(), "biz-ui", "superadmin", 1691, comment)
	if err != nil {
		log.Fatalf("Error creating comment: %v", err)
	}
}

func Test_CreateReview(t *testing.T) {
	if err := initClient(context.Background()); err != nil {
		color.Red("Failed to init GitHub client: %s", err)
		return
	}

	// List only your pending reviews on the pull request
	if err := deletePreviousPendingReview(context.Background(), "biz-ui", "superadmin", 1691); err != nil {
		color.Red("Failed to dismiss previous review: %s", err)
		return
	}
	draftComments := []*github.DraftReviewComment{
		//{
		//	Body: github.String("This is the drfat review comments 1."),
		//	Path: github.String("func/configuration/config/env.prd.us-east-1.yml"),
		//	Line: github.Int(20),
		//},
		{
			Body: github.String("-- 🤖 AI Generated --"),
			Path: github.String("func/configuration/config/env.prd.us-east-1.yml"),
			Line: github.Int(29),
		},
	}

	comment := &github.PullRequestReviewRequest{
		CommitID: github.String("91a7b754b0ebd920873edd7d6c40f2058290a388"),
		Comments: draftComments,
	}

	_, _, err := client.PullRequests.CreateReview(context.Background(), "biz-ui", "superadmin", 1691, comment)
	if err != nil {
		log.Fatalf("Error creating comment: %v", err)
	}

	_, _, err = client.PullRequests.CreateComment(context.Background(), "biz-ui", "superadmin", 1691, &github.PullRequestComment{
		//Body: github.String("-- 🤖 AI Generated --"),
		Path:     github.String("func/configuration/config/env.prd.us-east-1.yml"),
		Line:     github.Int(30),
		Body:     github.String("This is an extra comment."),
		CommitID: github.String("91a7b754b0ebd920873edd7d6c40f2058290a388"),
	})
	if err != nil {
		log.Fatalf("Error creating comment: %v", err)
	}
}

func Test_ListComments(t *testing.T) {
	if err := initClient(context.Background()); err != nil {
		color.Red("Failed to init GitHub client: %s", err)
		return
	}

	ctx := context.Background()
	owner := "biz-ui"
	repo := "superadmin"
	pullRequestNumber := 1691

	reviews, _, err := client.PullRequests.ListReviews(ctx, owner, repo, pullRequestNumber, nil)
	if err != nil {
		color.Red("Failed to list reviews: %s", err)
		return
	}

	// Print details of your pending reviews
	var pendingReview *github.PullRequestReview
	for _, review := range reviews {
		// Check if the review is pending (state "PENDING")
		if *review.State == "PENDING" {
			pendingReview = review
		}
	}

	if pendingReview == nil {
		color.Red("No pending review found")
		return
	}
	// List comments on the pending review
	comments, _, err := client.PullRequests.ListReviewComments(ctx, owner, repo, pullRequestNumber, *pendingReview.ID, nil)
	if err != nil {
		color.Red("Failed to list comments: %s", err)
		return
	}
	color.Green("Comments on the pending review: %v", comments)

}
