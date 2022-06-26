package main


func main() {

	var o = init_options()

	o.arg_analyze()

	o.env_analyze()

	o.execute()

}