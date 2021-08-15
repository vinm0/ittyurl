// Populate copyright year
const copyrightElement = document.getElementById("copyright")
let d = new Date();
copyrightElement.innerHTML = d.getFullYear()


function copyUrl() {
    /* Get the text field */
    const copyText = document.getElementById("copy-url");
    const copyBtn = document.getElementById("copy-btn");

    // Activate copy-url
    copyText.classList.remove("invisible");
    copyBtn.classList.remove("shake");

    /* Select the text field */
    copyText.select();
    copyText.setSelectionRange(0, 99999); /*For mobile devices*/

    /* Copy the text inside the text field */
    document.execCommand("copy");

    // Deactivate copy-url
    copyText.classList.add("invisible");

    /* Change copy button to read "Copied!" */
    copyBtn.innerHTML = "Copied!";
    copyBtn.classList.add("shake")
  };