package cmd

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(runCmd)

}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run dockerlayer",
	Long:  `Run dockerlayer`,
	Run: func(cmd *cobra.Command, args []string) {
		main2(args)
	},
}

func main2(args []string) {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	var isFiltering = false
	var filteringWord = ""
	var argc = len(args)
	if argc < 1 {
		fmt.Println("too few argments")
		os.Exit(2)
	} else if argc == 2 {
		isFiltering = true
		filteringWord = args[1]
	}
	fmt.Println("Image is [" + args[0] + "]")

	historyList, err := client.ImageHistory(args[0])
	if err != nil {
		fmt.Println("Error happens at client.ImageHistory")
		os.Exit(2)
	}
	// First CreatedBy newer to older
	var createdByHistoryList = []string{}
	//
	var etcpasswdmap = map[string]string{}
	var etcgroupmap = map[string]string{}
	// WORKAROUND
	var etcpasswdstring = ""
	var etcgroupstring = ""
	// WORKAROUND
	var etcpasswduidnamemap = map[int]string{}
	var etcgroupgidnamemap = map[int]string{}
	// layer ID -> json
	var jsonMap = map[string]string{}
	var allJSONMap = map[string]string{}
	var layerTarMap = map[string][]*tar.Header{}
	var first = false
	for _, history := range historyList {
		createdByHistoryList = append(createdByHistoryList, history.CreatedBy)
		image, err := client.InspectImage(history.ID)
		if err != nil || first {
			// continue to create CreatedBy list
			continue
		}
		if image != nil {
			var buf bytes.Buffer
			opts := docker.ExportImageOptions{Name: image.ID, OutputStream: &buf}
			err := client.ExportImage(opts)
			if err != nil {
				fmt.Println("Error happens at client.ExportImage")
				os.Exit(2)
			}
			r := bytes.NewReader(buf.Bytes())
			tr := tar.NewReader(r)
			var header *tar.Header
			for {
				header, err = tr.Next()
				if err == io.EOF {
					// end of image tar
					break
				}
				if err != nil {
					fmt.Println(err)
					fmt.Println("Error at tar extract")
					os.Exit(2)
				}

				buf2 := new(bytes.Buffer)
				if _, err = io.Copy(buf2, tr); err != nil {
					fmt.Println(err)
				}
				//
				var field = strings.Split(header.Name, "/")
				var layerID = field[0]
				// find json and layer.tar
				if strings.HasSuffix(header.Name, "/json") {
					var jsonstring = buf2.String()
					//var imagestring = ""
					allJSONMap[layerID] = jsonstring
					if strings.Index(jsonstring, "\"Image\":\"\"") != -1 || !first {
						jsonMap[layerID] = jsonstring
					}
				} else if strings.HasSuffix(header.Name, "layer.tar") {

					r2 := bytes.NewReader(buf2.Bytes())
					layerTar := tar.NewReader(r2)
					for {
						layerTarHeader, e4 := layerTar.Next()
						if e4 == io.EOF {
							// end of layer.tar
							break
						}
						layerTarMap[layerID] = append(layerTarMap[layerID], layerTarHeader)
						// read for etcpasswd and etcgroup
						if layerTarHeader.Name == "etc/passwd" || layerTarHeader.Name == "etc/group" {
							layerTarBuffer := new(bytes.Buffer)
							if _, err = io.Copy(layerTarBuffer, layerTar); err != nil {
								fmt.Println(err)
							}
							// TODO store json data
							var etcfilejsonstring = layerTarBuffer.String()
							// fmt.Println("FILENAME=[" + layerTarHeader.Name + "]")
							// fmt.Println(etcfilejsonstring)
							if layerTarHeader.Name == "etc/passwd" {
								etcpasswdmap[layerID] = etcfilejsonstring

								if etcpasswdstring == "" {
									etcpasswdstring = etcfilejsonstring
									var lines = strings.Split(etcpasswdstring, "\n")
									for _, line := range lines {
										var fields = strings.Split(line, ":")
										if len(fields) > 2 {
											var uid, _ = strconv.Atoi(fields[2])
											etcpasswduidnamemap[uid] = fields[0]
										}
									}
								}
							} else if layerTarHeader.Name == "etc/group" {
								etcgroupmap[layerID] = etcfilejsonstring
								if etcgroupstring == "" {
									etcgroupstring = etcfilejsonstring
									var lines = strings.Split(etcgroupstring, "\n")
									for _, line := range lines {
										var fields = strings.Split(line, ":")
										if len(fields) > 2 {
											var gid, _ = strconv.Atoi(fields[2])
											etcgroupgidnamemap[gid] = fields[0]
										}
									}
								}
							}
						}
					}
				}
			}
		}
		first = true

	}
	// fmt.Println("---- createdBy")
	// for _, createdBy := range createdByHistoryList {
	// 	fmt.Println(createdBy)
	// }
	// construct parent map
	//fmt.Println("---- jsonMap")
	var relation = map[string]string{}
	//var key string
	var value string
	var noparentid string
	for _, value = range jsonMap {
		//fmt.Println(key + " -> " + value)
		//dec := json.NewDecoder(v)
		var f interface{}
		//dec.Decode(&d)
		json.Unmarshal([]byte(value), &f)
		//fmt.Printf("%+v\n", d)
		m := f.(map[string]interface{})
		var id = m["id"]
		var parent = m["parent"]
		if parent == nil {
			noparentid = id.(string)
		} else {
			var p = parent.(string)
			_, ok := relation[p]
			if !ok {
				relation[p] = id.(string)
			}
		}
	}

	// fmt.Println("---- tarinfo")
	// fmt.Println("Start from:" + noparentid)
	var currentid = noparentid

	for {
		//fmt.Println(currentid)
		// find next
		_, ok := relation[currentid]
		if !ok {
			break
		}
		currentid = relation[currentid]
	}
	// File History List key is filename and value is layerid
	var fileHistoryListMap = map[string][]string{}
	var layerHistoryList = []string{}
	// fmt.Println("---- reconstruct")
	// fmt.Println("Start from:" + noparentid)
	currentid = noparentid
	for {
		//
		for _, layerTarHeader := range layerTarMap[currentid] {
			var filename = layerTarHeader.Name
			filename = strings.Replace(filename, ".wh.", "", -1)
			fileHistoryListMap[filename] = append(fileHistoryListMap[filename], currentid)
		}
		//
		layerHistoryList = append(layerHistoryList, currentid)
		_, ok := relation[currentid]
		if !ok {
			break
		}
		// fmt.Println(currentid + "->" + relation[currentid])
		currentid = relation[currentid]

	}
	//

	//
	// fmt.Println("---- layer")
	// reverse
	// SliceTricks · golang/go Wiki
	// https://github.com/golang/go/wiki/SliceTricks
	for i := len(layerHistoryList)/2 - 1; i >= 0; i-- {
		opp := len(layerHistoryList) - 1 - i
		layerHistoryList[i], layerHistoryList[opp] = layerHistoryList[opp], layerHistoryList[i]
	}
	// for _, layer := range layerHistoryList {
	// 	fmt.Println(layer)
	// }
	// fmt.Println("---- createdBy layer")
	// for _, createdBy := range createdByHistoryList {
	// 	fmt.Println(createdBy)
	// }
	var createdByListMap = map[string][]string{}
	var currentLayerID = ""
	var currentLayerIndex = 0
	for _, createdBy := range createdByHistoryList {
		if strings.Index(createdBy, "#(nop)") != -1 && strings.Index(createdBy, "ADD") == -1 && strings.Index(createdBy, "COPY") == -1 {
			//fmt.Println(strings.Repeat("-", 12) + " " + createdBy)
		} else {
			currentLayerID = layerHistoryList[currentLayerIndex]
			currentLayerIndex++
		}
		createdByListMap[currentLayerID] = append(createdByListMap[currentLayerID], createdBy)
	}
	//
	var layerIndex = 0

	for _, createdBy := range createdByHistoryList {
		if strings.Index(createdBy, "#(nop)") != -1 && strings.Index(createdBy, "ADD") == -1 && strings.Index(createdBy, "COPY") == -1 {
			//fmt.Println(strings.Repeat("-", 12) + " " + createdBy)
		} else {
			var layerID = layerHistoryList[layerIndex]
			for _, savedCreatedBy := range createdByListMap[layerID] {
				fmt.Println(layerID[:12] + "  " + savedCreatedBy)
			}
			//fmt.Println(allJSONMap[layerID])

			for _, layerTarHeader := range layerTarMap[layerID] {
				var filename = layerTarHeader.Name
				var deleteflag = false
				if strings.Index(filename, ".wh.") != -1 {
					deleteflag = true
					filename = strings.Replace(filename, ".wh.", "", -1)
				}
				//
				var status = "A"

				// check add or changes
				var fileHistoryList = fileHistoryListMap[filename]
				var fileHistoryIndex = 0
				var fileSize int64
				fileSize = 0
				for _, fileHistoryID := range fileHistoryList {
					var lTH2Size int64
					lTH2Size = 0
					for _, lTH2 := range layerTarMap[fileHistoryID] {
						if lTH2.Name == layerTarHeader.Name || lTH2.Name == filename {
							lTH2Size = lTH2.Size
							break
						}
					}
					fileSize = lTH2Size - fileSize
					if fileHistoryID == layerID {
						if fileHistoryIndex == 0 {

						} else {
							status = "C"

						}
						break
					}
					fileHistoryIndex++
				}

				if deleteflag {
					status = "D"
				}
				// calc size
				//
				var isOutput = false
				if isFiltering {
					if m, _ := regexp.MatchString(filteringWord, filename); m {
						isOutput = true
					}
				} else {
					isOutput = true
				}
				if isOutput {
					//fmt.Println(status + " " + filename + " " + strconv.FormatInt(fileSize, 10) + " " + strconv.Itoa(layerTarHeader.Uid) + "(" + layerTarHeader.Uname + ")" + ":" + strconv.Itoa(layerTarHeader.Gid) + "(" + layerTarHeader.Gname + ")" + " " + strconv.FormatInt(layerTarHeader.Mode, 8))
					fmt.Println(status + " " + filename + " " + strconv.FormatInt(fileSize, 10) + " " + strconv.Itoa(layerTarHeader.Uid) + "(" + etcpasswduidnamemap[layerTarHeader.Uid] + ")" + ":" + strconv.Itoa(layerTarHeader.Gid) + "(" + etcgroupgidnamemap[layerTarHeader.Gid] + ")" + " " + strconv.FormatInt(layerTarHeader.Mode, 8))
				}
			}

			layerIndex++
		}
	}

}
