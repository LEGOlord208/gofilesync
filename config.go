package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/legolord208/gofilesync/api"
	"github.com/legolord208/stdutil"
)

const port = ":8752"

var website = template.Must(template.New("gofilesync").Parse(html))

func initWebserver() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/schedule":
			str := trim(r.PostFormValue("min"))
			if checkmissing(w, "min", str) {
				return
			}
			min, err := strconv.Atoi(str)
			if err != nil {
				write(w, "That's not a number")
				return
			}

			if min < 0 {
				write(w, "Number is negative")
			}
			if min > 10000 {
				write(w, "Number is more than 10000")
			}

			data.Schedule = min
			schedule(min)
			write(w, saveData())
		case "/add":
			src := trim(r.PostFormValue("src"))
			if checkmissing(w, "src", src) {
				return
			}
			dst := trim(r.PostFormValue("dst"))
			if checkmissing(w, "dst", dst) {
				return
			}

			data.Locations = append(data.Locations, location{Src: src, Dst: dst})
			go func() {
				err := gofilesync.ForceSync(src, dst)
				if err != nil {
					status(true, err.Error())
				}
			}()
			write(w, saveData())
		case "/remove":
			str := trim(r.PostFormValue("index"))
			if checkmissing(w, "index", str) {
				return
			}
			index, err := strconv.Atoi(str)
			if err != nil {
				write(w, "That's not a number")
				return
			}

			if index < 0 || index >= len(data.Locations) {
				write(w, "Index not within bounds")
				return
			}

			data.Locations = append(data.Locations[:index], data.Locations[index+1:]...)
			write(w, saveData())
		case "/force-sync":
			str := trim(r.PostFormValue("index"))
			if checkmissing(w, "index", str) {
				return
			}
			index, err := strconv.Atoi(str)
			if err != nil {
				write(w, "That's not a number")
				return
			}

			if index < 0 || index >= len(data.Locations) {
				write(w, "Index not within bounds")
				return
			}

			loc := data.Locations[index]

			go func() {
				err := gofilesync.ForceSync(loc.Src, loc.Dst)
				if err != nil {
					status(true, err.Error())
				}
			}()
			write(w, "Started force sync...")
		default:
			website.Execute(w, data)
		}
	})
	err := http.ListenAndServe(port, nil)
	if err != nil {
		stdutil.PrintErr("Could not serve website", err)
	}
}

func write(w io.Writer, str string) {
	w.Write([]byte(str + "\n"))
}

func trim(str string) string {
	return strings.TrimSpace(str)
}
func checkmissing(w io.Writer, key, value string) bool {
	if value == "" {
		write(w, "'"+key+"' missing.")
		return true
	}
	return false
}

const html = `
<!DOCTYPE html>
<html>
	<head>
		<title>gofilesync</title>
		<style>
body {
	margin: 0;
	padding: 0;
	font-size: 20px
}
header {
	background: #3B83FF;
	color: white;
	text-align: center;
}
header > div:first-child {
	font-size: 40px;
}

.card {
	position: absolute;
	border: 1px solid black;
	border-radius: 10px;
	padding: 10px;
}
#card-schedule {
	left: 5%;
	top: 30%
}
#card-location-add {
	right: 25%;
	top: 20%
}
#card-locations {
	right: 20%;
	bottom: 20%
}
		</style>
	</head>
	<body>
		<header>
			<div>gofilesync</div>
			<div>Control Panel</div>
		</header>
		<div class="card" id="card-schedule">
			<h2>Backup Schedule:</h2>

			<p>Backup every <input type="text" style="width: 30px;text-align: center;" placeholder="X" value="{{.Schedule}}" /> minutes</p>
			<button>Save</button>
		</div>
		<div class="card" id="card-location-add">
			<h2>Backup Location Add</h2>

			<label>Source: <input type="text" id="src" placeholder="C:\A\File" /></label><br />
			<label>Destination: <input type="text" id="dst" placeholder="C:\A\File" /></label><br />
			<button>Save</button>
		</div>
		<div class="card" id="card-locations">
			<h2>Backup Locations</h2>

			<select>
				{{range $i, $loc := .Locations}}
					<option value="{{$i}}">{{$loc.Src}} ({{$loc.Dst}})</option>
				{{end}}
			</select>
			<button id="force-sync">Force sync</button>
			<button id="remove">Remove</button>
		</div>
		<script>
function makexhttp(url, reload) {
	var xhttp = new XMLHttpRequest()
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4) {
			alert(this.responseText);
			if (reload) {
				location.reload()
			}
		}
	}
	xhttp.open("POST", url)
	xhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded")
	return xhttp
}

document.querySelector("#card-schedule > button").addEventListener("click", function() {
	var val = document.querySelector("#card-schedule input").value;
	console.log(val);

	var xhttp = makexhttp("schedule");
	xhttp.send("min=" + encodeURIComponent(val));
})
document.querySelector("#card-location-add > button").addEventListener("click", function() {
	var val1 = document.getElementById("src").value;
	var val2 = document.getElementById("dst").value;
	console.log(val1 + ", " + val2);

	var xhttp = makexhttp("add", true);
	xhttp.send("src=" + encodeURIComponent(val1) + "&dst=" + encodeURIComponent(val2));
})
document.getElementById("force-sync").addEventListener("click", function() {
	var val = document.querySelector("#card-locations > select").value;
	console.log(val);

	var xhttp = makexhttp("force-sync");
	xhttp.send("index=" + encodeURIComponent(val));
})
document.getElementById("remove").addEventListener("click", function() {
	var val = document.querySelector("#card-locations > select").value;
	console.log(val);

	var xhttp = makexhttp("remove", true);
	xhttp.send("index=" + encodeURIComponent(val));
})
		</script>
	</body>
</html>
`
