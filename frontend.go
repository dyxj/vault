package main

//<button type="submit" formaction="/decrypt">Decrypt</button>

const defPageConst string = `<form action="/encrypt" method="post" enctype="multipart/form-data">
upload a file to encrypt<br>
<input type="file" name="usrfile"><br>
<input type="text" name="cryptKey1" value=""><br>
<input type="text" name="cryptKey2" value=""><br>
<button type="submit">Encrypt</button>
</form>
<br>
<form action="/decrypt" method="post" enctype="multipart/form-data">
upload a file to decrypt<br>
<input type="file" name="usrfile"><br>
<input type="text" name="cryptKey1" value="">
<br>
<button type="submit">decrypt</button><br>
</form>
`
