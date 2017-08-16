package main

import (
	"cloudify"
	"cloudify/utils"
	"flag"
	"fmt"
	"os"
)

var host string
var user string
var password string
var tenant string

func basicOptions(name string) *flag.FlagSet {
	var commonFlagSet *flag.FlagSet
	commonFlagSet = flag.NewFlagSet(name, flag.ExitOnError)
	commonFlagSet.StringVar(&host, "host", "localhost", "Manager host name")
	commonFlagSet.StringVar(&user, "user", "admin", "Manager user name")
	commonFlagSet.StringVar(&password, "password", "secret", "Manager user password")
	commonFlagSet.StringVar(&tenant, "tenant", "default_tenant", "Manager tenant")
	return commonFlagSet
}

func infoOptions(args, options []string) int {
	defaultError := "state/version subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "state":
		{
			operFlagSet := basicOptions("status state")
			operFlagSet.Parse(options)
			stat := cloudify.GetStatus(host, user, password, tenant)

			fmt.Printf("Retrieving manager services status... [ip=%v]\n", host)
			fmt.Printf("Manager status: %v\n", stat.Status)
			fmt.Printf("Services:\n")
			var lines [][]string = make([][]string, len(stat.Services))
			for pos, service := range stat.Services {
				lines[pos] = make([]string, 2)
				lines[pos][0] = service.DisplayName
				lines[pos][1] = service.Status()
			}
			utils.PrintTable([]string{"service", "status"}, lines)
		}
	case "version":
		{
			operFlagSet := basicOptions("status version")
			operFlagSet.Parse(options)

			ver := cloudify.GetVersion(host, user, password, tenant)
			fmt.Printf("Retrieving manager services version... [ip=%v]\n", host)
			utils.PrintTable([]string{"Version", "Edition"}, [][]string{{ver.Version, ver.Edition}})
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}

func blueprintsOptions(args, options []string) int {
	defaultError := "list/delete subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}
	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("blueprints list")
			operFlagSet.Parse(options)
			blueprints := cloudify.GetBlueprints(host, user, password, tenant)
			var lines [][]string = make([][]string, len(blueprints.Items))
			for pos, blueprint := range blueprints.Items {
				lines[pos] = make([]string, 7)
				lines[pos][0] = blueprint.Id
				lines[pos][1] = blueprint.Description
				lines[pos][2] = blueprint.MainFileName
				lines[pos][3] = blueprint.CreatedAt
				lines[pos][4] = blueprint.UpdatedAt
				lines[pos][5] = blueprint.Tenant
				lines[pos][6] = blueprint.CreatedBy
			}
			utils.PrintTable([]string{"id", "description", "main_file_name", "created_at", "updated_at", "tenant_name", "created_by"}, lines)
		}
	case "delete":
		{
			operFlagSet := basicOptions("blueprints delete")
			if len(args) < 4 {
				fmt.Println("Blueprint Id requered")
				return 1
			}

			operFlagSet.Parse(options)
			blueprint := cloudify.DeleteBlueprints(host, user, password, tenant, args[3])
			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 7)
			lines[0][0] = blueprint.Id
			lines[0][1] = blueprint.Description
			lines[0][2] = blueprint.MainFileName
			lines[0][3] = blueprint.CreatedAt
			lines[0][4] = blueprint.UpdatedAt
			lines[0][5] = blueprint.Tenant
			lines[0][6] = blueprint.CreatedBy
			utils.PrintTable([]string{"id", "description", "main_file_name", "created_at", "updated_at", "tenant_name", "created_by"}, lines)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}

func deploymentsOptions(args, options []string) int {
	defaultError := "list/create/delete subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("deployments list")
			operFlagSet.Parse(options)
			deployments := cloudify.GetDeployments(host, user, password, tenant)
			var lines [][]string = make([][]string, len(deployments.Items))
			for pos, deployment := range deployments.Items {
				lines[pos] = make([]string, 6)
				lines[pos][0] = deployment.Id
				lines[pos][1] = deployment.BlueprintId
				lines[pos][2] = deployment.CreatedAt
				lines[pos][3] = deployment.UpdatedAt
				lines[pos][4] = deployment.Tenant
				lines[pos][5] = deployment.CreatedBy
			}
			utils.PrintTable([]string{"id", "blueprint_id", "created_at", "updated_at", "tenant_name", "created_by"}, lines)
		}
	case "create":
		{
			operFlagSet := basicOptions("deployments list <deployment id>")
			if len(args) < 4 {
				fmt.Println("Deployment Id requered")
				return 1
			}

			var blueprint string
			operFlagSet.StringVar(&blueprint, "blueprint", "", "The unique identifier for the blueprint")

			operFlagSet.Parse(options)

			var depl cloudify.CloudifyDeploymentPost
			depl.BlueprintId = blueprint
			depl.Inputs = map[string]interface{}{}
			deployment := cloudify.CreateDeployments(host, user, password, tenant, args[3], depl)

			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 6)
			lines[0][0] = deployment.Id
			lines[0][1] = deployment.BlueprintId
			lines[0][2] = deployment.CreatedAt
			lines[0][3] = deployment.UpdatedAt
			lines[0][4] = deployment.Tenant
			lines[0][5] = deployment.CreatedBy
			utils.PrintTable([]string{"id", "blueprint_id", "created_at", "updated_at", "tenant_name", "created_by"}, lines)
		}
	case "delete":
		{
			operFlagSet := basicOptions("deployments delete <deployment id>")
			if len(args) < 4 {
				fmt.Println("Deployment Id requered")
				return 1
			}

			operFlagSet.Parse(options)
			deployment := cloudify.DeleteDeployments(host, user, password, tenant, args[3])
			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 6)
			lines[0][0] = deployment.Id
			lines[0][1] = deployment.BlueprintId
			lines[0][2] = deployment.CreatedAt
			lines[0][3] = deployment.UpdatedAt
			lines[0][4] = deployment.Tenant
			lines[0][5] = deployment.CreatedBy
			utils.PrintTable([]string{"id", "blueprint_id", "created_at", "updated_at", "tenant_name", "created_by"}, lines)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}

func executionsOptions(args, options []string) int {
	defaultError := "list/start subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("executions list")

			var deployment string
			operFlagSet.StringVar(&deployment, "deployment", "", "The unique identifier for the deployment")
			operFlagSet.Parse(options)

			var options = map[string]string{}
			if deployment != "" {
				options["deployment_id"] = deployment
			}
			executions := cloudify.GetExecutions(host, user, password, tenant, options)
			var lines [][]string = make([][]string, len(executions.Items))
			for pos, execution := range executions.Items {
				lines[pos] = make([]string, 8)
				lines[pos][0] = execution.Id
				lines[pos][1] = execution.WorkflowId
				lines[pos][2] = execution.Status
				lines[pos][3] = execution.DeploymentId
				lines[pos][4] = execution.CreatedAt
				lines[pos][5] = execution.Error
				lines[pos][6] = execution.Tenant
				lines[pos][7] = execution.CreatedBy
			}
			utils.PrintTable([]string{"id", "workflow_id", "status", "deployment_id", "created_at", "error", "tenant_name", "created_by"}, lines)
		}
	case "start":
		{
			operFlagSet := basicOptions("executions start <workflow id>")
			if len(args) < 4 {
				fmt.Println("Workflow Id requered")
				return 1
			}

			var deployment string
			operFlagSet.StringVar(&deployment, "deployment", "", "The unique identifier for the deployment")
			operFlagSet.Parse(options)

			var exec cloudify.CloudifyExecutionPost
			exec.WorkflowId = args[3]
			exec.DeploymentId = deployment

			execution := cloudify.PostExecution(host, user, password, tenant, exec)

			var lines [][]string = make([][]string, 1)
			lines[0] = make([]string, 8)
			lines[0][0] = execution.Id
			lines[0][1] = execution.WorkflowId
			lines[0][2] = execution.Status
			lines[0][3] = execution.DeploymentId
			lines[0][4] = execution.CreatedAt
			lines[0][5] = execution.Error
			lines[0][6] = execution.Tenant
			lines[0][7] = execution.CreatedBy
			utils.PrintTable([]string{"id", "workflow_id", "status", "deployment_id", "created_at", "error", "tenant_name", "created_by"}, lines)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}

func eventsOptions(args, options []string) int {
	defaultError := "list subcommand is required"

	if len(args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	switch args[2] {
	case "list":
		{
			operFlagSet := basicOptions("events list")
			var blueprint string
			var deployment string
			var execution string
			operFlagSet.StringVar(&blueprint, "blueprint", "", "The unique identifier for the blueprint")
			operFlagSet.StringVar(&deployment, "deployment", "", "The unique identifier for the deployment")
			operFlagSet.StringVar(&execution, "execution", "", "The unique identifier for the execution")
			operFlagSet.Parse(options)

			var options = map[string]string{}
			if blueprint != "" {
				options["blueprint_id"] = blueprint
			}
			if deployment != "" {
				options["deployment_id"] = deployment
			}
			if execution != "" {
				options["execution_id"] = execution
			}
			if blueprint != "" {
				options["blueprint_id"] = blueprint
			}
			events := cloudify.GetEvents(host, user, password, tenant, options)
			var lines [][]string = make([][]string, len(events.Items))
			for pos, event := range events.Items {
				lines[pos] = make([]string, 5)
				lines[pos][0] = event.Timestamp
				lines[pos][1] = event.DeploymentId
				lines[pos][2] = event.NodeInstanceId
				lines[pos][3] = event.Operation
				lines[pos][4] = event.Message
			}
			utils.PrintTable([]string{"Timestamp", "Deployment", "InstanceId", "Operation", "Message"}, lines)
		}
	default:
		{
			fmt.Println(defaultError)
			return 1
		}
	}
	return 0
}

func main() {

	args, options := utils.CliArgumentsList(os.Args)
	defaultError := "Supported only: status, version, blueprints, deployments, executions, events"
	if len(args) < 2 {
		fmt.Println(defaultError)
		return
	}

	switch args[1] {
	case "status":
		{
			os.Exit(infoOptions(args, options))
		}
	case "blueprints":
		{
			os.Exit(blueprintsOptions(args, options))
		}
	case "deployments":
		{
			os.Exit(deploymentsOptions(args, options))
		}
	case "executions":
		{
			os.Exit(executionsOptions(args, options))
		}
	case "events":
		{
			os.Exit(eventsOptions(args, options))
		}
	default:
		{
			fmt.Println(defaultError)
			os.Exit(1)
		}
	}
}
