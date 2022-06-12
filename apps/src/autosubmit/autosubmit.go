package main

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

var (
	jobInfo Job
	yamlPath string
	repoAddr string
	jobId string
	resultPath string
)


type Job struct {
	Suite string `yaml:"suite"`
	Category string `yaml:"category"`
	RepoAddr string `yaml:"repo_addr"`
	CustomRepoName string `yaml:"custom_repo_name"`
	DockerImage string `yaml:"docker_image"`
	Arch string `yaml:"arch"`
	MountRepoAddr []string `yaml:"mount_repo_addr"`
	MountRepoName []string`yaml:"mount_repo_name"`
	TestBox string `yaml:"testbox"`
}


func init(){
	args := os.Args
	if len(args) != 2{
		fmt.Println("第一个参数: yaml文件位置")
		// fmt.Println("第二个参数: repo_addr")
		os.Exit(1)
	}else {
		yamlPath = args[1]
		// repoAddr = args[2]
	}
}


func updateYaml(){
	fileData,err := ioutil.ReadFile(yamlPath)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(fileData, &jobInfo)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	jobInfo.RepoAddr = repoAddr
	fileData,err = yaml.Marshal(&jobInfo)
	if err != nil{
		fmt.Printf("RepoAddr写入失败: %s\n",err.Error())
		os.Exit(1)
	}
	contents := string(fileData)
	lines := strings.Split(contents, "\n")
	rpmBuild := "rpmbuild: "
	if len(lines) < 2{
		fmt.Println("请检查yaml格式")
		os.Exit(1)
	}
	lines = append(lines[:2],append([]string{rpmBuild},lines[2:]...)...)
	var newFileData string
	for _,line := range lines{
		newFileData = fmt.Sprintf("%s%s\n",newFileData,line)
	}
	err = ioutil.WriteFile(yamlPath,[]byte(newFileData),0666)
	if err != nil{
		fmt.Printf("文件写入失败: %s\n",err.Error())
		os.Exit(1)
	}
}

func submit(){
	cmd := exec.Command("submit", yamlPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if !strings.Contains( string(out),"got job id="){
		fmt.Println("跳转失败")
		fmt.Println("输入字符串需要包含got job id=zx.xxxxxxxx")
		os.Exit(1)
	}
	jobId = strings.Split(string(out),"got job id=")[1]
	if len(jobId) != 0{
		fmt.Printf("submit successful id=%s\n",jobId)
	}
}


func debugInfo(){
	repoItems := strings.Split(repoAddr,"/")
	repoInfo := repoItems[len(repoItems)-1]
	fmt.Println(fmt.Sprintf("    %s\n    id=%s\n----------------------------------------------\n",repoInfo,jobId))
}


func createLog (){
	yamlPathItems := strings.Split(yamlPath,"/")
	logName := fmt.Sprintf("%s.log",yamlPathItems[len(yamlPathItems)-1])
	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("日志更新失败", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	write := bufio.NewWriter(file)
	repoItems := strings.Split(repoAddr,"/")
	repoInfo := repoItems[len(repoItems)-1]
	repoLog := fmt.Sprintf("    %s\n    id=%s\n----------------------------------------------\n",repoInfo,jobId)
	_, err = write.WriteString(repoLog)
	if err != nil {
		fmt.Println("日志更新错误",err)
		return
	}
	err = write.Flush()
	if err != nil {
		fmt.Println("日志写入错误",err)
		return
	}
}

func getRepoAddr(addrSlice []string) (err error){
	if len(addrSlice) <=2{
		repoAddr = addrSlice[0]
		err = nil
	}else{
		err = errors.New("仓库格式错误，不应该携带一个以上的'&'")
	}
	return
}
func main() {
	canExit := true
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, os.Interrupt)
	waitExit(exitSignal,canExit)
	for{
		fmt.Printf("请输入 repo_addr: ")
		_, err := fmt.Scan(&repoAddr)
		aliyunRepo := strings.Split(repoAddr,"&")
		err = getRepoAddr(aliyunRepo)
		for err != nil{
			fmt.Printf("输入错误:%s\n",err.Error())
			fmt.Printf("请重新输入 repo_addr: ")
			_, err = fmt.Scan(&repoAddr)
			aliyunRepo = strings.Split(repoAddr,"&")
			err = getRepoAddr(aliyunRepo)
		}
		canExit = false
		updateYaml()
		submit()
		debugInfo()
		createLog()
		canExit = true
	}
}

func waitExit(exit chan os.Signal,canExit bool) {
	go func() {
		<-exit
		for canExit{
			os.Exit(0)
		}
	}()
}