packagemain

import(
	"net/http"
)

func helloGoHandler(w http.RespontseWriter,r *http.Request){
	fmt.Fprintf(w,"HelloGo....")
}

func hiGoHandler(w http.RespontseWriter,r *http.Request){
	fmt.Fprintf(w,"HiGo....")
}

func main(){
	http.HandleFunc("/hellogo",helloGoHandler)
	http.HandleFunc("/higo",hiGoHandler)

	fmt.Pritf("Server start at port 3030 ...")

	if err := http.ListenAndServe(":3030",nil);err != nil {
		fmt.Printf("Serversart error: %v",err)
	}
}