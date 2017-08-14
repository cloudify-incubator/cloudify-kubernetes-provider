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
	commonFlagSet = flag.NewFlagSet("name", flag.ExitOnError)
	commonFlagSet.StringVar(&host, "host", "localhost", "Manager host name")
	commonFlagSet.StringVar(&user, "user", "admin", "Manager user name")
	commonFlagSet.StringVar(&password, "password", "secret", "Manager user password")
	commonFlagSet.StringVar(&tenant, "tenant", "default_tenant", "Manager tenant")
	return commonFlagSet
}

func infoOptions() int {
	defaultError := "state/version subcommand is required"

	if len(os.Args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	operFlagSet := basicOptions("status")

	operFlagSet.Parse(os.Args[3:])

	switch os.Args[2] {
	case "state":
		{
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

func blueprintsOptions() int {
	defaultError := "list/delete subcommand is required"

	if len(os.Args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	operFlagSet := basicOptions("blueprints")

	switch os.Args[2] {
	case "list":
		{
			operFlagSet.Parse(os.Args[3:])
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
			if len(os.Args) < 4 {
				fmt.Println("Blueprint Id requered")
				return 1
			}

			operFlagSet.Parse(os.Args[4:])
			blueprint := cloudify.DeleteBlueprints(host, user, password, tenant, os.Args[3])
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

func deploymentsOptions() int {
	defaultError := "list/create/delete subcommand is required"

	if len(os.Args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	operFlagSet := basicOptions("deployments")

	switch os.Args[2] {
	case "list":
		{
			operFlagSet.Parse(os.Args[3:])
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
			if len(os.Args) < 4 {
				fmt.Println("Deployment Id requered")
				return 1
			}

			var blueprint string
			operFlagSet.StringVar(&blueprint, "blueprint", "", "The unique identifier for the blueprint")

			operFlagSet.Parse(os.Args[4:])

			var depl cloudify.CloudifyDeploymentPost
			depl.BlueprintId = blueprint
			depl.Inputs = map[string]string{}
			deployment := cloudify.CreateDeployments(host, user, password, tenant, os.Args[3], depl)

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
			if len(os.Args) < 4 {
				fmt.Println("Deployment Id requered")
				return 1
			}

			operFlagSet.Parse(os.Args[4:])
			deployment := cloudify.DeleteDeployments(host, user, password, tenant, os.Args[3])
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

func executionsOptions() int {
	defaultError := "list/start subcommand is required"

	if len(os.Args) < 3 {
		fmt.Println(defaultError)
		return 1
	}

	operFlagSet := basicOptions("executions")

	switch os.Args[2] {
	case "list":
		{
			operFlagSet.Parse(os.Args[3:])
			executions := cloudify.GetExecutions(host, user, password, tenant)
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

			if len(os.Args) < 4 {
				fmt.Println("Workflow Id requered")
				return 1
			}

			var deployment string
			operFlagSet.StringVar(&deployment, "deployment", "", "The unique identifier for the deployment")
			operFlagSet.Parse(os.Args[4:])

			var exec cloudify.CloudifyExecutionPost
			exec.WorkflowId = os.Args[3]
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

func main() {
	defaultError := "Supported only: status, version, blueprints, deployments, executions, executions-install"
	if len(os.Args) < 2 {
		fmt.Println(defaultError)
		return
	}

	switch os.Args[1] {
	case "status":
		{
			os.Exit(infoOptions())
		}
	case "blueprints":
		{
			os.Exit(blueprintsOptions())
		}
	case "deployments":
		{
			os.Exit(deploymentsOptions())
		}
	case "executions":
		{
			os.Exit(executionsOptions())
		}
	default:
		{
			fmt.Println(defaultError)
			os.Exit(1)
		}
	}
}
