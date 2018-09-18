function search() {
	term = document.getElementById("search").value;
	bucket = document.getElementById("search").getAttribute("name");

	url = "/search/" + bucket + "/" + term;
	open(url, "_self");
}
