package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alajmo/mani/core"
	"github.com/alajmo/mani/core/dao"
)

func runCmd(config *dao.Config, configErr *error) *cobra.Command {
	var runFlags core.RunFlags

	cmd := cobra.Command{
		Use:   "run <task> [flags]",
		Short: "Run tasks",
		Long: `Run tasks.

The tasks are specified in a mani.yaml file along with the projects you can target.`,

		Example: `  # Run task 'pwd' for all projects
  mani run pwd --all-projects

  # Checkout branch 'development' for all projects that have tag 'backend'
  mani run checkout -t backend branch=development`,

		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			core.CheckIfError(*configErr)
			run(args, config, &runFlags)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if *configErr != nil {
				return []string{}, cobra.ShellCompDirectiveDefault
			}

			return config.GetTaskNames(), cobra.ShellCompDirectiveNoFileComp
		},
	}

	cmd.Flags().BoolVar(&runFlags.Describe, "describe", false, "Print task information")
	cmd.Flags().BoolVar(&runFlags.DryRun, "dry-run", false, "don't execute any task, just print the output of the task to see what will be executed")
	cmd.Flags().BoolVarP(&runFlags.Edit, "edit", "e", false, "Edit task")
	cmd.Flags().BoolVar(&runFlags.Parallel, "parallel", false, "Run tasks in parallel for each project")

	cmd.Flags().StringVarP(&runFlags.Output, "output", "o", "", "Output text|table|html|markdown")
	err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if *configErr != nil {
			return []string{}, cobra.ShellCompDirectiveDefault
		}

		valid := []string{"text", "table", "html", "markdown"}
		return valid, cobra.ShellCompDirectiveDefault
	})
	core.CheckIfError(err)

	cmd.Flags().BoolVarP(&runFlags.Cwd, "cwd", "k", false, "current working directory")

	cmd.Flags().BoolVar(&runFlags.AllProjects, "all-projects", false, "target all projects")
	cmd.Flags().BoolVar(&runFlags.AllDirs, "all-dirs", false, "target all dirs")

	cmd.Flags().StringSliceVarP(&runFlags.Projects, "projects", "p", []string{}, "target projects by their name")
	err = cmd.RegisterFlagCompletionFunc("projects", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if *configErr != nil {
			return []string{}, cobra.ShellCompDirectiveDefault
		}

		projects := config.GetProjectNames()
		return projects, cobra.ShellCompDirectiveDefault
	})
	core.CheckIfError(err)

	cmd.Flags().StringSliceVarP(&runFlags.Dirs, "dirs", "d", []string{}, "target directories by their name")
	err = cmd.RegisterFlagCompletionFunc("dirs", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if *configErr != nil {
			return []string{}, cobra.ShellCompDirectiveDefault
		}

		dirs := config.GetDirNames()
		return dirs, cobra.ShellCompDirectiveDefault
	})
	core.CheckIfError(err)

	cmd.Flags().StringSliceVarP(&runFlags.Paths, "paths", "g", []string{}, "target directories by their path")
	err = cmd.RegisterFlagCompletionFunc("paths", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if *configErr != nil {
			return []string{}, cobra.ShellCompDirectiveDefault
		}

		projectPaths := config.GetProjectPaths()
		dirPaths := config.GetDirPaths()

		options := append(projectPaths, dirPaths...)

		return options, cobra.ShellCompDirectiveDefault
	})
	core.CheckIfError(err)

	cmd.Flags().StringSliceVarP(&runFlags.Tags, "tags", "t", []string{}, "target entities by their tag")
	err = cmd.RegisterFlagCompletionFunc("tags", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if *configErr != nil {
			return []string{}, cobra.ShellCompDirectiveDefault
		}

		tags := config.GetTags()
		return tags, cobra.ShellCompDirectiveDefault
	})
	core.CheckIfError(err)

	return &cmd
}

func run(
	args []string,
	config *dao.Config,
	runFlags *core.RunFlags,
) {
	var taskNames []string
	var userArgs []string
	// Seperate user arguments from task names
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			userArgs = append(userArgs, arg)
		} else {
			taskNames = append(taskNames, arg)
		}
	}

	if runFlags.Edit {
		if len(args) > 0 {
			config.EditTask(taskNames[0])
			return
		} else {
			config.EditTask("")
			return
		}
	}

	for _, name := range taskNames {
		task, err := config.GetTask(name)
		core.CheckIfError(err)

		projectEntities, dirEntities := config.GetTaskEntities(task, *runFlags)

		if len(projectEntities) == 0 && len(dirEntities) == 0 {
			fmt.Println("No targets")
		} else {
			if len(projectEntities) > 0 {
				entityList := dao.EntityList{
					Type:     "Project",
					Entities: projectEntities,
				}

				task.RunTask(entityList, userArgs, config, runFlags)
			}

			if len(dirEntities) > 0 {
				entityList := dao.EntityList{
					Type:     "Directory",
					Entities: dirEntities,
				}
				task.RunTask(entityList, userArgs, config, runFlags)
			}
		}
	}
}
