package cmd

//
//var logsCmd = &cobra.Command{
//	Use:   "logs",
//	Short: "Logs",
//	Args:  cobra.ExactArgs(1),
//	RunE: func(cmd *cobra.Command, args []string) error {
//		LoginGate()
//		ctx := context.Background()
//		client := config.GetAPIClient()
//
//		projectPath := args[0]
//		project, err := client.GetProjectByPath(ctx, projectPath)
//		if err != nil {
//			return err
//		}
//
//		logs, err := client.GetProjectLogs(ctx, project.ID.String())
//		if err != nil {
//			return err
//		}
//
//		for _, log := range logs {
//			fmt.Println(log)
//		}
//
//		return nil
//	},
//}
//
//func init() {
//	rootCmd.AddCommand(logsCmd)
//}
