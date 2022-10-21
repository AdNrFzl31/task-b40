let number = 0;

function addBlog() {
  let name = document.getElementById("input-name").value;
  let description = document.getElementById("input-description").value;

  let startDate = new Date(document.getElementById("input-startDate").value);
  let endDate = new Date(document.getElementById("input-endDate").value);

  let nodeJs = document.getElementById("input-nodeJs").checked;
  let javascript = document.getElementById("input-javascript").checked;
  let reactJs = document.getElementById("input-reactJs").checked;
  let html5 = document.getElementById("input-html5").checked;

  if (nodeJs) {
    nodeJs = "fa-brands fa-node-js fa-xl";
  } else {
    nodeJs = "";
  }

  if (javascript) {
    javascript = "fa-brands fa-square-js fa-xl";
  } else {
    javascript = "";
  }

  if (reactJs) {
    reactJs = "fa-brands fa-react fa-xl";
  } else {
    reactJs = "";
  }

  if (html5) {
    html5 = "fa-brands fa-html5 fa-xl";
  } else {
    html5 = "";
  }

  let image = document.getElementById("input-image").files;

  if (image.length == 0) {
    return alert("Gambar harus di isi!!");
  } else {
    image = URL.createObjectURL(image[0]);
  }
  console.log(image);

  document.getElementById("content").innerHTML += `
    <div class="card" id="card${number}">
      <a href="blog-detail.html">
        <img
          src="${image}"
          alt="gambar avatar"
          class="image"
        />
        <div class="card-body">
          <div>
            <h3>${name}</h3>
            <span> durasi : ${getDistanceTime(startDate, endDate)} </span>
          </div>
          <p>
            ${description}
          </p>
          <div class="icon-group">
            ${nodeJs && `<i class='icon ${nodeJs}'></i>`}
            ${javascript && `<i class='icon ${javascript}'></i>`}
            ${reactJs && `<i class='icon ${reactJs}'></i>`}
            ${html5 && `<i class='icon ${html5}'></i>`}
          </div>
        </div>
      </a>
      <div class="button-group">
        <button type="button" class="btn-card">edit</button>
        <button type="button" class="btn-card" onclick="removeButton('card${number}')">delete</button>
    </div>
  </div>
    `;

    number++
}



function getDistanceTime(starDate, endDate) {
  let timeNow = new Date(endDate);
  let timePost = new Date(starDate);


  let distance = timeNow - timePost; 
  console.log(distance);

  let milisecond = 1000;
  let secondInHours = 3600; 
  let hoursInDay = 24; 
  let dayInMonth = 30;
  let monthInYear = 12;  

  let distanceYear = Math.floor(
    distance / (milisecond * secondInHours * hoursInDay * dayInMonth * monthInYear)
  );

  let distanceMonth = Math.floor(distance / (milisecond * 60 * 60 * 24 * 12));
  let distanceDay = Math.floor(distance / (milisecond * 60 * 60 * 24));
  let distanceHours = Math.floor(distance / (milisecond * 60 * 60));
  let distanceMinutes = Math.floor(distance / (milisecond * 60));
  let distanceSecond = Math.floor(distance / milisecond);

  if (distanceYear > 0) {
    return `${distanceYear} Year ago`  
  } else if (distanceMonth > 0) {
    return `${distanceMonth} Month ago`
  } else if (distanceDay > 0) {
    return `${distanceDay} Day ago`;
  } else if (distanceHours > 0) {
    return `${distanceHours} Hours ago`;
  } else if (distanceMinutes > 0) {
    return `${distanceMinutes} Minutes ago`;
  } else {
    return `${distanceSecond} Seconds ago`;
  }
}


function removeButton(card) {
  let element = document.getElementById(card);
  element.remove();
}
