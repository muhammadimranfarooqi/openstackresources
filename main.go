package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/term"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
//	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
//	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/quotas"
//	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/quotas"
        "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
//	"github.com/gophercloud/gophercloud/openstack/clustering/v1/nodes"
//        "github.com/gophercloud/gophercloud/openstack/clustering/v1/policies"

)

func main() {
        var verbose bool
        flag.BoolVar(&verbose, "verbose", false, "verbose output")
        var name string
        flag.StringVar(&name, "name", "", "user or app name")
        var password string
        flag.StringVar(&password, "password", "", "user password or app secret")
        var endpoint string
        flag.StringVar(&endpoint, "endpoint", "https://keystone.cern.ch/v3", "endpoint")
        var project string
        flag.StringVar(&project, "project", "CMS Webtools Mig", "project")
        flag.Parse()
        log.SetFlags(0)
        log.SetFlags(log.Lshortfile)
        if name == "" {
                var err error
                name, password, err = credentials()
                if err != nil {
                        log.Fatal("unable to read credentials")
                }
        }
        run(endpoint, name, password, project, verbose)

}

// helper function to get user credentials from stdin
func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Println("")

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

// helper function to run our workflow on openstack
func run(endpoint, username, password, project string, verbose bool) {


	var opts gophercloud.AuthOptions
	var err error
	scope := &gophercloud.AuthScope{
		ProjectName: project,
		DomainName:  "default",
		DomainID:    "default",
	}
	opts = gophercloud.AuthOptions{
			IdentityEndpoint: endpoint,
			Username:         username,
			Password:         password,
			DomainID:         "default",
			Scope:            scope,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Fatal("auth client failure: ", err)
	}
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		log.Fatal("compute client error", err)
	}

	// list existing servers
	pager := servers.List(client, servers.ListOpts{})
	pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)
		if err != nil {
			log.Println("extract servers: ", err)
			return false, err
		}
		for _, s := range serverList {
			if verbose {
				fmt.Println(s.ID, s.Name,  s.Status)
			}
		}
		return true, nil
	})
	// Block storage Quota
/*	quotaset, err := quotasets.GetUsage(client, os.Getenv("OS_PROJECT_ID")).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset)
*/
	//quotaset2, err := quotasets.Get(client, "tenant-id").Extract()
	quotaset2, err := quotasets.Get(client, os.Getenv("OS_PROJECT_ID")).Extract()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", quotaset2)

/*
        quotasInfo, err := quotas.Get(client, os.Getenv("OS_PROJECT_ID")).Extract()
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("quotas: %#v\n", quotasInfo)
*/


/*
listOpts := nodes.ListOpts{
	Name: "cmsweb-test6",
}

allPages, err := nodes.List(client, listOpts).AllPages()
if err != nil {
	panic(err)
}

allNodes, err := nodes.ExtractNodes(allPages)
if err != nil {
	panic(err)
}

for _, node := range allNodes {
	fmt.Printf("%+v\n", node)
}

*/

/*
listOpts := policies.ListOpts{
	Limit: 2,
}

allPages, err := policies.List(client, listOpts).AllPages()
if err != nil {
	panic(err)
}

allPolicies, err := policies.ExtractPolicies(allPages)
if err != nil {
	panic(err)
}

for _, policy := range allPolicies {
	fmt.Printf("%+v\n", policy)
}
*/
}
