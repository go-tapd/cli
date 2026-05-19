package cmd

import (
	"errors"

	"github.com/go-tapd/cli/internal/app"
	"github.com/go-tapd/tapd"
	"github.com/spf13/cobra"
)

func newBoardCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "board",
		Short: "Work with TAPD boards",
	}

	cmd.AddCommand(
		newBoardCardCmd(rt),
		newBoardColumnsCmd(rt),
	)
	return cmd
}

func newBoardCardCmd(rt *app.Runtime) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "card",
		Short: "Work with board cards",
	}

	cmd.AddCommand(
		newBoardCardCreateCmd(rt),
		newBoardCardListCmd(rt),
		newBoardCardUpdateCmd(rt),
	)
	return cmd
}

func newBoardCardCreateCmd(rt *app.Runtime) *cobra.Command {
	var flags boardCardMutationFlags

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a board card",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.CreateBoardCardRequest{
				WorkspaceID: tapd.Ptr(flags.workspaceID),
				BoardID:     tapd.Ptr(flags.boardID),
				ColumnID:    tapd.Ptr(flags.columnID),
				Name:        tapd.Ptr(flags.name),
			}
			applyBoardCardCreateFlags(request, flags)

			card, _, err := client.BoardService.CreateBoardCard(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, boardCardHeaders(), boardCardRows([]*tapd.BoardCard{card}), card)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	cmd.Flags().Int64Var(&flags.boardID, "board-id", 0, "board ID")
	cmd.Flags().Int64Var(&flags.columnID, "column-id", 0, "column ID")
	cmd.Flags().StringVar(&flags.name, "name", "", "card name")
	addBoardCardMutationFlags(cmd, &flags, false)
	_ = cmd.MarkFlagRequired("board-id")
	_ = cmd.MarkFlagRequired("column-id")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newBoardCardListCmd(rt *app.Runtime) *cobra.Command {
	var flags boardCardQueryFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List board cards",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request, err := newGetBoardCardsRequest(flags)
			if err != nil {
				return err
			}

			cards, _, err := client.BoardService.GetBoardCards(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, boardCardHeaders(), boardCardRows(cards), cards)
		},
	}

	addBoardCardQueryFlags(cmd, &flags)
	return cmd
}

func newBoardCardUpdateCmd(rt *app.Runtime) *cobra.Command {
	var flags boardCardMutationFlags

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a board card",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseInt64Arg("board card id", args[0])
			if err != nil {
				return err
			}
			if !hasBoardCardUpdate(flags) {
				return errors.New("at least one board card field flag is required")
			}

			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.UpdateBoardCardRequest{
				ID:          tapd.Ptr(id),
				WorkspaceID: tapd.Ptr(flags.workspaceID),
			}
			applyBoardCardUpdateFlags(request, flags)

			card, _, err := client.BoardService.UpdateBoardCard(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, boardCardHeaders(), boardCardRows([]*tapd.BoardCard{card}), card)
		},
	}

	newWorkspaceFlag(cmd, &flags.workspaceID)
	addBoardCardMutationFlags(cmd, &flags, true)
	return cmd
}

func newBoardColumnsCmd(rt *app.Runtime) *cobra.Command {
	var (
		workspaceID int
		ids         string
		name        string
		boardID     int64
		status      string
		created     string
		creator     string
		limit       int
		page        int
		fields      string
	)

	cmd := &cobra.Command{
		Use:   "columns",
		Short: "List board columns",
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, _, err := rt.NewClient()
			if err != nil {
				return err
			}

			request := &tapd.GetBoardColumnsRequest{
				WorkspaceID: tapd.Ptr(workspaceID),
				Limit:       tapd.Ptr(limit),
				Page:        tapd.Ptr(page),
				Fields:      fieldsMulti(fields),
			}
			if ids != "" {
				request.ID, err = strictInt64Multi("column IDs", ids)
				if err != nil {
					return err
				}
			}
			if name != "" {
				request.Name = tapd.Ptr(name)
			}
			if boardID > 0 {
				request.BoardID = tapd.Ptr(boardID)
			}
			if status != "" {
				request.Status = tapd.Ptr(status)
			}
			if created != "" {
				request.Created = tapd.Ptr(created)
			}
			if creator != "" {
				request.Creator = tapd.Ptr(creator)
			}

			columns, _, err := client.BoardService.GetBoardColumns(cmd.Context(), request)
			if err != nil {
				return err
			}

			return writeOutput(cmd, rt.OutputFormat, boardColumnHeaders(), boardColumnRows(columns), columns)
		},
	}

	newWorkspaceFlag(cmd, &workspaceID)
	newListFlags(cmd, &limit, &page)
	cmd.Flags().StringVar(&ids, "ids", "", "comma separated column IDs")
	cmd.Flags().StringVar(&name, "name", "", "filter by column name")
	cmd.Flags().Int64Var(&boardID, "board-id", 0, "board ID")
	cmd.Flags().StringVar(&status, "status", "", "filter by status")
	cmd.Flags().StringVar(&created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&creator, "creator", "", "filter by creator")
	cmd.Flags().StringVar(&fields, "fields", "", "comma separated fields")
	return cmd
}

type boardCardMutationFlags struct {
	workspaceID int
	boardID     int64
	columnID    int64
	name        string
	owner       string
	cc          string
	status      string
	begin       string
	due         string
	label       int64
	description string
}

func addBoardCardMutationFlags(cmd *cobra.Command, flags *boardCardMutationFlags, update bool) {
	if update {
		cmd.Flags().Int64Var(&flags.boardID, "board-id", 0, "board ID")
		cmd.Flags().Int64Var(&flags.columnID, "column-id", 0, "column ID")
		cmd.Flags().StringVar(&flags.name, "name", "", "card name")
	}
	cmd.Flags().StringVar(&flags.owner, "owner", "", "card owner")
	cmd.Flags().StringVar(&flags.cc, "cc", "", "card participants")
	cmd.Flags().StringVar(&flags.status, "status", "", "card status")
	cmd.Flags().StringVar(&flags.begin, "begin", "", "begin date")
	cmd.Flags().StringVar(&flags.due, "due", "", "due date")
	cmd.Flags().Int64Var(&flags.label, "label", 0, "label ID")
	cmd.Flags().StringVar(&flags.description, "description", "", "card description")
}

func applyBoardCardCreateFlags(request *tapd.CreateBoardCardRequest, flags boardCardMutationFlags) {
	if flags.owner != "" {
		request.Owner = tapd.Ptr(flags.owner)
	}
	if flags.cc != "" {
		request.CC = tapd.Ptr(flags.cc)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	if flags.begin != "" {
		request.Begin = tapd.Ptr(flags.begin)
	}
	if flags.due != "" {
		request.Due = tapd.Ptr(flags.due)
	}
	if flags.label > 0 {
		request.Label = tapd.Ptr(flags.label)
	}
	if flags.description != "" {
		request.Description = tapd.Ptr(flags.description)
	}
}

func applyBoardCardUpdateFlags(request *tapd.UpdateBoardCardRequest, flags boardCardMutationFlags) {
	if flags.boardID > 0 {
		request.BoardID = tapd.Ptr(flags.boardID)
	}
	if flags.columnID > 0 {
		request.ColumnID = tapd.Ptr(flags.columnID)
	}
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.owner != "" {
		request.Owner = tapd.Ptr(flags.owner)
	}
	if flags.cc != "" {
		request.CC = tapd.Ptr(flags.cc)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	if flags.begin != "" {
		request.Begin = tapd.Ptr(flags.begin)
	}
	if flags.due != "" {
		request.Due = tapd.Ptr(flags.due)
	}
	if flags.label > 0 {
		request.Label = tapd.Ptr(flags.label)
	}
	if flags.description != "" {
		request.Description = tapd.Ptr(flags.description)
	}
}

func hasBoardCardUpdate(flags boardCardMutationFlags) bool {
	return flags.boardID > 0 ||
		flags.columnID > 0 ||
		flags.name != "" ||
		flags.owner != "" ||
		flags.cc != "" ||
		flags.status != "" ||
		flags.begin != "" ||
		flags.due != "" ||
		flags.label > 0 ||
		flags.description != ""
}

type boardCardQueryFlags struct {
	workspaceID int
	limit       int
	page        int
	fields      string
	ids         string
	boardID     int64
	columnID    int64
	owner       string
	cc          string
	status      string
	name        string
	created     string
	begin       string
	due         string
	label       int64
}

func addBoardCardQueryFlags(cmd *cobra.Command, flags *boardCardQueryFlags) {
	newWorkspaceFlag(cmd, &flags.workspaceID)
	newListFlags(cmd, &flags.limit, &flags.page)
	cmd.Flags().StringVar(&flags.fields, "fields", "", "comma separated fields")
	cmd.Flags().StringVar(&flags.ids, "ids", "", "comma separated card IDs")
	cmd.Flags().Int64Var(&flags.boardID, "board-id", 0, "board ID")
	cmd.Flags().Int64Var(&flags.columnID, "column-id", 0, "column ID")
	cmd.Flags().StringVar(&flags.owner, "owner", "", "filter by owner")
	cmd.Flags().StringVar(&flags.cc, "cc", "", "filter by participants")
	cmd.Flags().StringVar(&flags.status, "status", "", "filter by status")
	cmd.Flags().StringVar(&flags.name, "name", "", "filter by card name")
	cmd.Flags().StringVar(&flags.created, "created", "", "filter by created time expression")
	cmd.Flags().StringVar(&flags.begin, "begin", "", "filter by begin date expression")
	cmd.Flags().StringVar(&flags.due, "due", "", "filter by due date expression")
	cmd.Flags().Int64Var(&flags.label, "label", 0, "filter by label ID")
}

func newGetBoardCardsRequest(flags boardCardQueryFlags) (*tapd.GetBoardCardsRequest, error) {
	request := &tapd.GetBoardCardsRequest{
		WorkspaceID: tapd.Ptr(flags.workspaceID),
		Limit:       tapd.Ptr(flags.limit),
		Page:        tapd.Ptr(flags.page),
		Fields:      fieldsMulti(flags.fields),
	}
	if flags.ids != "" {
		ids, err := strictInt64Multi("card IDs", flags.ids)
		if err != nil {
			return nil, err
		}
		request.ID = ids
	}
	if flags.boardID > 0 {
		request.BoardID = tapd.Ptr(flags.boardID)
	}
	if flags.columnID > 0 {
		request.ColumnID = tapd.Ptr(flags.columnID)
	}
	if flags.owner != "" {
		request.Owner = tapd.Ptr(flags.owner)
	}
	if flags.cc != "" {
		request.CC = tapd.Ptr(flags.cc)
	}
	if flags.status != "" {
		request.Status = tapd.Ptr(flags.status)
	}
	if flags.name != "" {
		request.Name = tapd.Ptr(flags.name)
	}
	if flags.created != "" {
		request.Created = tapd.Ptr(flags.created)
	}
	if flags.begin != "" {
		request.Begin = tapd.Ptr(flags.begin)
	}
	if flags.due != "" {
		request.Due = tapd.Ptr(flags.due)
	}
	if flags.label > 0 {
		request.Label = tapd.Ptr(flags.label)
	}
	return request, nil
}

func boardCardHeaders() []string {
	return []string{"ID", "Name", "BoardID", "ColumnID", "Status", "Owner", "Due"}
}

func boardCardRows(cards []*tapd.BoardCard) [][]string {
	rows := make([][]string, 0, len(cards))
	for _, item := range cards {
		rows = append(rows, []string{
			item.ID,
			item.Name,
			item.BoardID,
			item.ColumnID,
			item.Status,
			stringValue(item.Owner),
			stringValue(item.Due),
		})
	}
	return rows
}

func boardColumnHeaders() []string {
	return []string{"ID", "Name", "BoardID", "Status", "Creator", "Created"}
}

func boardColumnRows(columns []*tapd.BoardColumn) [][]string {
	rows := make([][]string, 0, len(columns))
	for _, item := range columns {
		rows = append(rows, []string{
			item.ID,
			item.Name,
			item.BoardID,
			item.Status,
			item.Creator,
			item.Created,
		})
	}
	return rows
}
