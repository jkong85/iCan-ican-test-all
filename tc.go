package main

import (
	"log"
	"os/exec"
	"strings"
	"strconv"
)

func main() {
	del(2)
	config(2)
	log.Println("<<<<<<<<<<<<<<<<<<<<<========================>>>>>>>>>>>>>>>>")
	show(2)
}

func config(num int){
	log.Println("Start to config the TC policy on All Nodes!")
	for i := 0; i<num; i++{
		node_name := "ican-" + strconv.Itoa(i)
		log.Println("========================")
		log.Println("Config the node: " + node_name)
		// for each node, we config the VM's interface, Br interface and all the container internace and Veth		
		node_cmd_prefix := "ssh " + node_name + " sudo "

		// config the VM's interface
		configEth(node_cmd_prefix, "enp0s9")
	 	// config the Br
		configEth(node_cmd_prefix, "br")

		// config the Docker's interface and Veth
		cmd := node_cmd_prefix + " docker ps -q"
		ids := exe_cmd_full(cmd)
		if ids == ""{
			log.Println("No pod on node " + node_name)
			continue
		}
		for _, container_id := range strings.Split(ids, "\n") {
			log.Println("container_id is: " + container_id)
			if container_id == "" {
				log.Println("container_id is nil")
				continue
			}
			//get container pid
			cmd_docker := node_cmd_prefix + " docker inspect -f {{.State.Pid}} " + container_id
			container_pid := strings.Trim(exe_cmd_full(cmd_docker), "\n")
			log.Println("container pid is: " + container_pid)
			//for each container, set its interface's TC policy
			// tc qdisc add dev eth0 root tbf rate 100mbit latency 50ms burst 100k
			cmd_docker = node_cmd_prefix + " nsenter -t " + container_pid + " -n " +  " tc qdisc add dev eth0 root tbf rate 100mbit latency 50ms burst 100k"
			log.Println(exe_cmd_full(cmd_docker))
		}
		// then we find out the veth and config the tc policy
		cmd = node_cmd_prefix + " ifconfig | grep veth "

		for _, line := range strings.Split(exe_cmd_full(cmd), "\n") {
			if line == ""{
				log.Println("veth is nil")		
				continue
			}
			veth_name := strings.Split(line, " ")[0]
			log.Println("veth name is:" + veth_name)
			cmd_veth := node_cmd_prefix + " tc qdisc add dev " + veth_name + " root tbf rate 100mbit latency 50ms burst 100k"

			_ = exe_cmd_full(cmd_veth)
		}

	}
}

func configEth(node_cmd_prefix string, dev string) {
	/*
	   tc qdisc add dev eth0 root handle 1: htb default 12

	   tc class add dev eth0 parent 1: classid 1:1 htb rate 100kbps ceil 100kbps

	   tc class add dev eth0 parent 1:1 classid 1:10 htb rate 30kbps ceil 100kbps 
	   tc class add dev eth0 parent 1:1 classid 1:11 htb rate 10kbps ceil 100kbps 
	   tc class add dev eth0 parent 1:1 classid 1:12 htb rate 60kbps ceil 100kbps 

	   tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip src 172.17.0.2 flowid 1:10
	   tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip src 173.17.0.2 match ip dport 80 0xffff flowid 1:11
	   tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip src 1.2.3.4 flowid 1:12
	 */
	cmd := node_cmd_prefix + " tc qdisc add dev  " + dev + " root handle 1: htb default 12 "
	exe_cmd_full(cmd)
	cmd = node_cmd_prefix + " tc class add dev  " + dev + " parent 1: classid 1:1 htb rate 100kbps ceil 100kbps "
	exe_cmd_full(cmd)
	cmd = node_cmd_prefix + " tc class add dev  " + dev + " parent 1:1 classid 1:10 htb rate 30kbps ceil 100kbps "
	exe_cmd_full(cmd)
	cmd = node_cmd_prefix + " tc class add dev  " + dev + " parent 1:1 classid 1:11 htb rate 10kbps ceil 100kbps "
	exe_cmd_full(cmd)
	cmd = node_cmd_prefix + " tc class add dev  " + dev + " parent 1:1 classid 1:12 htb rate 60kbps ceil 100kbps "
	exe_cmd_full(cmd)
	cmd = node_cmd_prefix + " tc filter add dev " + dev + " protocol ip parent 1:0 prio 1 u32 match ip src 172.17.0.2 flowid 1:10 "
	exe_cmd_full(cmd)
	cmd = node_cmd_prefix + " tc filter add dev " + dev + " protocol ip parent 1:0 prio 1 u32 match ip src 173.17.0.2 match ip dport 80 0xffff flowid 1:11"
	exe_cmd_full(cmd)
	cmd = node_cmd_prefix + " tc filter add dev " + dev + " protocol ip parent 1:0 prio 1 u32 match ip src 1.2.3.4 flowid 1:12 "
	exe_cmd_full(cmd)
}
func showEth(node_cmd_prefix string, dev string){

	cmd := node_cmd_prefix + " tc qdisc show dev " + dev
	out := exe_cmd_full(cmd)
	log.Println(out)

	cmd = node_cmd_prefix + " tc class show dev " + dev
	out = exe_cmd_full(cmd)
	log.Println(out)

	cmd = node_cmd_prefix + " tc filter show dev " + dev
	out = exe_cmd_full(cmd)
	log.Println(out)
}



func show(num int){
	for i := 0; i<num; i++{
		node_name := "ican-" + strconv.Itoa(i)
		log.Println("========================")
		log.Println("Show the config of the node: " + node_name)
		// for each node, we config the VM's interface, Br interface and all the container internace and Veth		
		node_cmd_prefix := "ssh " + node_name + " sudo "
		// show the configuration for VM's interface
		showEth(node_cmd_prefix, "enp0s9")
		showEth(node_cmd_prefix, "br")

		cmd := node_cmd_prefix + " docker ps -q"
		ids := exe_cmd_full(cmd)
		if ids == ""{
			log.Println("No pod on node " + node_name)
			continue
		}
		for _, container_id := range strings.Split(ids, "\n") {
			log.Println("container_id is: " + container_id)
			if container_id == "" {
				log.Println("container_id is nil")
				continue
			}
			//get container pid
			cmd_docker := node_cmd_prefix + " docker inspect -f {{.State.Pid}} " + container_id
			container_pid := strings.Trim(exe_cmd_full(cmd_docker), "\n")
			log.Println("container pid is: " + container_pid)
			//for each container, set its interface's TC policy
			cmd_docker = node_cmd_prefix + " nsenter -t " + container_pid + " -n " +  " tc qdisc show dev eth0 "
			log.Println(exe_cmd_full(cmd_docker))
		}
		// then we find out the veth and config the tc policy
		cmd = node_cmd_prefix + " ifconfig | grep veth "

		for _, line := range strings.Split(exe_cmd_full(cmd), "\n") {
			if line == ""{
				log.Println("veth is nil")
				continue
			}
			veth_name := strings.Split(line, " ")[0]
			log.Println("veth name is:" + veth_name)
			cmd_veth := node_cmd_prefix + " tc qdisc show dev " + veth_name
			_ = exe_cmd_full(cmd_veth)
		}

	}
}


func del(num int){
	for i := 0; i<num; i++{
		node_name := "ican-" + strconv.Itoa(i)
		log.Println("========================")
		log.Println("delete the configuration of the node: " + node_name)
		// for each node, we config the VM's interface, Br interface and all the container internace and Veth		
		node_cmd_prefix := "ssh " + node_name + " sudo "
		// show the configuration for VM's interface

		cmd := node_cmd_prefix + " tc qdisc del dev enp0s9 root "
		exe_cmd_full(cmd)

		cmd = node_cmd_prefix + " tc qdisc del dev br root "
		exe_cmd_full(cmd)

		cmd = node_cmd_prefix + " docker ps -q"
		ids := exe_cmd_full(cmd)
		if ids == ""{
			log.Println("No pod on node " + node_name)
			continue
		}
		for _, container_id := range strings.Split(ids, "\n") {
			log.Println("container_id is: " + container_id)
			if container_id == "" {
				log.Println("container_id is nil")
				continue
			}
			//get container pid
			cmd_docker := node_cmd_prefix + " docker inspect -f {{.State.Pid}} " + container_id
			container_pid := strings.Trim(exe_cmd_full(cmd_docker), "\n")
			log.Println("container pid is: " + container_pid)
			//for each container, set its interface's TC policy
			cmd_docker = node_cmd_prefix + " nsenter -t " + container_pid + " -n " +  " tc qdisc del dev eth0 root"
			log.Println(exe_cmd_full(cmd_docker))
		}
		// then we find out the veth and config the tc policy
		cmd = node_cmd_prefix + " ifconfig | grep veth "

		for _, line := range strings.Split(exe_cmd_full(cmd), "\n") {
			if line == ""{
				log.Println("veth is nil")		
				continue
			}
			veth_name := strings.Split(line, " ")[0]
			log.Println("veth name is:" + veth_name)
			cmd_veth := node_cmd_prefix + " tc qdisc del dev " + veth_name + " root "
			_ = exe_cmd_full(cmd_veth)
		}

	}
}



func exe_cmd_full(cmd string) string {
	log.Println("command is : ", cmd)
	//out, _ := exec.Command("sh", "-c", cmd).Output()
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("Error to exec CMD", cmd)
	}
	//log.Println("Output of command:", string(out))
	return string(out)
}

