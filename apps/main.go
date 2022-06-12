package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	jobInfo Job
	yamlPath string
	repoAddr string
	jobId string
	resultPath string
)
const (
	pathPrefix = "/srv/result/rpmbuild"
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
	if len(args) != 3{
		fmt.Println("第一个参数: yaml文件位置")
		fmt.Println("第二个参数: repo_addr")
		os.Exit(1)
	}else {
		yamlPath = args[1]
		repoAddr = args[2]
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
	fmt.Printf("---------submit---------\n %s ---------submit---------\n", string(out))
	jobId = strings.Split(string(out),"got job id=")[1]
}


func getResultPath(){
	timeStamp := time.Now().Format("2006-01-02")
	dockerImage := strings.Replace(jobInfo.DockerImage,":","-",-1)
	dockerImage = fmt.Sprintf("%s-%s",dockerImage,jobInfo.Arch)
	resultPath = fmt.Sprintf("%s/%s/%s/%s/%s/%s",pathPrefix,timeStamp,jobInfo.TestBox,dockerImage,jobInfo.Arch,jobId)
	fmt.Println("submit输出目录如下:")
	fmt.Println(resultPath)
}


func createLog (){
	yamlPathItems := strings.Split(yamlPath,"/")
	logName := fmt.Sprintf("%s.log",yamlPathItems[len(yamlPathItems)-1])
	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("日志更新失败", err)
		os.Exit(1)
	}
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	repoItems := strings.Split(repoAddr,"/")
	repoInfo := repoItems[len(repoItems)-1]
	repoLog := fmt.Sprintf("    %s\n    id=%s\n----------------------------------------------\n",repoInfo,jobId)
	write.WriteString(repoLog)
	write.Flush()
}

func main() {
	updateYaml()
	// submit()
	getResultPath()
	createLog()
}