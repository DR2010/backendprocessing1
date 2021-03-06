var prato

function danielonload() {
    document.getElementById("date").innerHTML = "Date: " + getdatetimer();

    showlinechart();
}

// --------------------------------------
//    Load Prices in Memory or Cache
// --------------------------------------
function loadPrices() {
    prato = new Array();
    prato[0] = { name: "Coxinha", price: "7" };
    prato[1] = { name: "BolodeCenoura", price: "3" };
    prato[2] = { name: "Refrigerante", price: "5" };
    prato[3] = { name: "BolodeAipim", price: "9" };
    prato[4] = { name: "PasteldeQueijo", price: "5" };
    prato[5] = { name: "Brigadeiro", price: "8" };
}

// --------------------------------------
//       Add New Item to Order
// --------------------------------------
function addNewItem() {

    loadPrices()

    var table = document.getElementById("myTable");
    var valueplate = document.getElementById("pratoname");
    var plateqtd = document.getElementById("pratoqtd");

    var lastRow = table.rows[table.rows.length];
    var lastRowNumber = table.rows.length;

    var row = table.insertRow(lastRow);

    var cell0 = row.insertCell(0);
    var cell1 = row.insertCell(1);
    var cell2 = row.insertCell(2);
    var cell3 = row.insertCell(3);

    x = '<input type=checkbox name=row' + lastRowNumber + ' id=checkitem' + lastRowNumber + ' value=' + valueplate + '>';

    var sel = document.getElementById("pratoname").selectedIndex;
    cell0.innerHTML = x;
    cell1.innerHTML = valueplate.nodeValue;
    cell1.innerHTML = valueplate.options[sel].text;
    cell2.innerHTML = plateqtd.value;
    cell3.innerHTML = plateqtd.value * prato[sel].price;
}

// --------------------------------------
//       Clean up fields
// --------------------------------------
function newOrder() {

    var orderID = document.getElementById("orderID");
    var orderClientName = document.getElementById("orderClientName");
    var orderDate = document.getElementById("orderDate");
    var orderTime = document.getElementById("orderTime");
    var eatmode = document.getElementById("EatMode");
    var status = document.getElementById("status");
    var message = document.getElementById("message");

    orderID.value = "";
    orderClientName.value = "";
    orderDate.value = getdate();
    orderTime.value = "";
    eatmode.value = "Eatin";
    status.value = "New Order";
    message.value = "Place new order";

}

// --------------------------------------
//   Show Selected Rows - Debug only
// --------------------------------------

function showSelectedRows() {

    var selchbox = []; // array that will store the value of selected checkboxes

    var table = document.getElementById("myTable");

    var lastRowNumber = table.rows.length - 1;

    for (var i = lastRowNumber; i >= 0; i--) {
        var chk = document.getElementById('checkitem' + i);

        var col2 = table.rows[i].cells[2].innerHTML;

        if (chk != null)
            if (chk.checked) alert(col2);

    }
    return selchbox;
}

// --------------------------------------
//   Remove Selected Rows from Order
// --------------------------------------

function removeSelectedRows() {
    // JavaScript & jQuery Course - http://coursesweb.net/javascript/
    var selchbox = []; // array that will store the value of selected checkboxes

    var table = document.getElementById("myTable");
    var lastRowNumber = table.rows.length;

    for (var i = lastRowNumber; i >= 0; i--) {
        var chk = document.getElementById('checkitem' + i);

        if (chk != null)
            if (chk.checked) table.deleteRow(i);

    }

    return selchbox;
}


// --------------------------------------
//          Save order JSON
// --------------------------------------
function saveOrder() {

    var orderID = document.getElementById("orderID");
    var orderClientName = document.getElementById("orderClientName");
    var orderDate = document.getElementById("orderDate");
    var orderTime = document.getElementById("orderTime");
    var eatmode = document.getElementById("EatMode");
    var status = document.getElementById("status");
    var message = document.getElementById("message");

    if (orderClientName.value == "") {
        message.value = "Order name is mandatory!"
        orderClientName.focus();
        return
    }


    if (orderID.value != "") {
        message.value = "Order already placed!"
        return
    }

    var oTable = document.getElementById('myTable');
    var rowLength = oTable.rows.length;

    if (rowLength == 1) {
        message.value = "Please add items!"
        return
    }


    var pratosselected = new Array();

    //loops through rows    
    for (i = 0; i < rowLength; i++) {

        var oCells = oTable.rows.item(i).cells;
        var cellLength = oCells.length;

        var pratoname = "";
        var quantidade = "";
        var preco = "";


        for (var j = 0; j < cellLength; j++) {

            var cellVal = oCells.item(j).innerHTML;
            if (j == 0) {
                action = cellVal;
            }
            if (j == 1) {
                pratoname = cellVal;
            }
            if (j == 2) {
                quantidade = cellVal;
            }
            if (j == 3) {
                preco = cellVal;

            }

            pratosselected[i] = { pratoname: pratoname, quantidade: quantidade, preco: preco };
        }
    }

    // Build the object - order
    // Post to the server or call web api
    var http = new XMLHttpRequest();
    var url = "/orderadd";

    status.value = "Placed";

    var paramsjson = JSON.stringify({
        orderID: orderID.value,
        orderClientName: orderClientName.value,
        orderDate: orderDate.value,
        orderTime: orderTime.value,
        eatmode: eatmode.value,
        status: status.value,
        pratos: pratosselected
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { //Call a function when the state changes.
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);
            message.value = "Order placed successfully."
            status.value = "Placed"

            var json_data = http.responseText;

            var contact = JSON.parse(json_data);
            orderID.value = contact.ID;

        }
    }

    http.send(paramsjson);

}

// --------------------------------------
//          Save order JSON
// --------------------------------------
function backToList() {

    window.location.replace("/orderlist");
}

function pad(num, size) {
    return ('000000000' + num).substr(-size);
}

function getdatetime() {
    var date = document.getElementById("orderDate");
    var time = document.getElementById("orderTime");

    var today = new Date();
    var dd = today.getDate();
    var mm = today.getMonth() + 1; //January is 0!
    var hh = pad(today.getHours(), 2);
    var min = pad(today.getMinutes(), 2);

    var yyyy = today.getFullYear();
    if (dd < 10) {
        dd = '0' + dd;
    }
    if (mm < 10) {
        mm = '0' + mm;
    }
    var today = yyyy + '-' + mm + '-' + dd;
    date.value = today;

    var hour = hh + ':' + min;
    time.value = hour;

}

function getdate() {

    var today = new Date();
    var dd = today.getDate();
    var mm = today.getMonth() + 1; //January is 0!
    var hh = pad(today.getHours(), 2);
    var min = pad(today.getMinutes(), 2);

    var yyyy = today.getFullYear();
    if (dd < 10) {
        dd = '0' + dd;
    }
    if (mm < 10) {
        mm = '0' + mm;
    }
    var today = yyyy + '-' + mm + '-' + dd;

    return today

}

function getdatetimer() {

    var today = new Date();
    var dd = today.getDate();
    var mm = today.getMonth() + 1; //January is 0!
    var hh = pad(today.getHours(), 2);
    var min = pad(today.getMinutes(), 2);

    var yyyy = today.getFullYear();
    if (dd < 10) {
        dd = '0' + dd;
    }
    if (mm < 10) {
        mm = '0' + mm;
    }
    var today = yyyy + '-' + mm + '-' + dd + ' ~ ' + hh + ':' + min;

    return today

}

function showchart() {
    var ctx = document.getElementById("myChart").getContext('2d');
    var myChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: ["Red", "Blue", "Yellow", "Green", "Purple", "Orange"],
            datasets: [{
                label: '# of Votes',
                data: [12, 19, 3, 5, 2, 3],
                backgroundColor: [
                    'rgba(255, 99, 132, 0.2)',
                    'rgba(54, 162, 235, 0.2)',
                    'rgba(255, 206, 86, 0.2)',
                    'rgba(75, 192, 192, 0.2)',
                    'rgba(153, 102, 255, 0.2)',
                    'rgba(255, 159, 64, 0.2)'
                ],
                borderColor: [
                    'rgba(255,99,132,1)',
                    'rgba(54, 162, 235, 1)',
                    'rgba(255, 206, 86, 1)',
                    'rgba(75, 192, 192, 1)',
                    'rgba(153, 102, 255, 1)',
                    'rgba(255, 159, 64, 1)'
                ],
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {
                        beginAtZero: true
                    }
                }]
            }
        }
    });
}

function showlinechart() {
    showanychart('line');
}

function showlinechartdate(){

    var fromDate = document.getElementById('fromDate');
    var toDate = document.getElementById('toDate');
    var currency = document.getElementById('currency');

    var url = "/btcmarketshistorylistdate?currency="+currency.value+"&fromDate="+fromDate.value+"&toDate="+toDate.value;

    window.location.href = url;
}

function showbarchart() {
    showanychart('bar');
}

function showanychart(chartype) {

    var theTbl = document.getElementById('tablecotacao');
    var arr = [];
    var labellist = [];

    // for (var i = 1; i < theTbl.rows.length; i++) {
    var st = theTbl.rows.length - 1;
    for (var i = st; i >= 1; i--) {
        // for (var i = 1; i < theTbl.rows.length; i++) {

        var value = theTbl.rows[i].cells[4].innerHTML;
        var valuecol0 = theTbl.rows[i].cells[0].innerHTML.substr(62, 3);

        if (valuecol0 == "ALL") {
            arr.push(Number(theTbl.rows[i].cells[3].innerHTML));
        } else {
            arr.push(Number(theTbl.rows[i].cells[2].innerHTML));
        }

        labellist.push(value.substr(32, 5));
    }

    showchartline(arr, labellist, chartype);
}

function showchartline(datalist, labellist, chartype) {

    var ctx = document.getElementById("myChart").getContext('2d');

    var myChart = new Chart(ctx, {
        type: chartype,
        data: {
            labels: labellist,
            datasets: [{
                label: 'Cotacao',
                data: datalist,
                backgroundColor: [
                    'rgba(255, 99, 132, 0.2)',
                    'rgba(54, 162, 235, 0.2)',
                    'rgba(255, 206, 86, 0.2)',
                    'rgba(75, 192, 192, 0.2)',
                    'rgba(153, 102, 255, 0.2)',
                    'rgba(255, 159, 64, 0.2)'
                ],
                borderColor: [
                    'rgba(255,99,132,1)',
                    'rgba(54, 162, 235, 1)',
                    'rgba(255, 206, 86, 1)',
                    'rgba(75, 192, 192, 1)',
                    'rgba(153, 102, 255, 1)',
                    'rgba(255, 159, 64, 1)'
                ],
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {}
                }]
            }
        }
    });

}



// --------------------------------------
//          Login
// --------------------------------------
function login() {

    var userid = document.getElementById("userid");
    var password = document.getElementById("password");
    var message = document.getElementById("message");

    if (userid.value == "") {
        userid.focus();
        return
    }

    if (password.value == "") {
        password.focus()
        return
    }

    // Build the object - order
    // Post to the server or call web api
    var http = new XMLHttpRequest();
    var url = "/login";

    var paramsjson = JSON.stringify({
        userid: userid.value,
        password: password.value
    });

    http.open("POST", url, true);

    //Send the proper header information along with the request
    http.setRequestHeader("Content-type", "application/json");

    http.onreadystatechange = function() { //Call a function when the state changes.
        if (http.readyState == 4 && http.status == 200) {
            console.log(http.responseText);
            message.value = "Order placed successfully."

            var json_data = http.responseText;
            var contact = JSON.parse(json_data);
        } else {
            message.value = "Invalid Credentials."
        }
    }

    http.send(paramsjson);

}