package main

const defPageConst string = `<form action="/encrypt" method="post" enctype="multipart/form-data">
upload a file<br>
<input type="file" name="usrfile"><br>
<input type="text" name="cryptKey1" value=""><br>
<input type="text" name="cryptKey2" value=""><br>
<button type="submit">Encrypt</button><br>
<button type="submit" formaction="/decrypt">Decrypt</button>
</form>
<br>
<br>
<h1></h1>`
