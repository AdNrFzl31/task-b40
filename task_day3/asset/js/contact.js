// materi

function submitData(e) {
  let name = document.getElementById("input-name").value;
  let email = document.getElementById("input-email").value;
  let phone = document.getElementById("input-phone").value;
  let subject = document.getElementById("input-subject").value;
  let message = document.getElementById("input-message").value;

  console.log(name);
  console.log(email);
  console.log(phone);
  console.log(subject);
  console.log(message);
  console.log(this);

  let emailReceiver = "adenurfaizal31@gmail.com";
  location.href =
    "mailto:" +
    emailReceiver +
    "?subject=" +
    subject +
    "&body= Hallo nama saya " +
    name +
    ", %0D%0A " +
    message +
    ", %0D%0A silahkan kontak ke nomor " +
    phone;
}
