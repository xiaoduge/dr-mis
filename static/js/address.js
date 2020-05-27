

window.onload = initProvince;

function initProvince(){
    $("#provinces").empty();
    var strSpace = '<option value="space"></option>';
    $("#provinces").append(strSpace);
    $.ajax({
        type: "GET",
        url: "/provincesinfo",
        contentType: 'application/json;charset=utf-8', //设置请求头信息
        dataType: "json",
        success: function (result, status) {
            console.log("status: " + status)
            var infos = result["regioninfo"];

            for(let i = 0; i < infos.length; i++){
                var str = '<option value="' + infos[i]["region_code"] + '">' + infos[i]["region_fullname"] +'</option>'
                $("#provinces").append(str);
            }
        },
        error: function (result) {
            console.log(result.msg);
        }
    });
};

function show_city(opt){
    $("#citys").empty();
    var strSpace = '<option value="space"></option>';
    $("#citys").append(strSpace);
    var value = opt.value;
    console.log("value: " + value) 
    $.ajax({
        type: "GET",
        url: "/citysinfo?belongs=" + value,
        contentType: 'application/json;charset=utf-8', //设置请求头信息
        dataType: "json",
        success: function (result, status) {
            console.log("status: " + status)
            var infos = result["regioninfo"];
            for(let i = 0; i < infos.length; i++){
                var str = '<option value="' + infos[i]["region_code"] + '">' + infos[i]["region_fullname"] +'</option>'
                $("#citys").append(str);
            }
            
        },
        error: function (result) {
            console.log(result.msg);
        }
    });
};

function show_county(opt){
    $("#countys").empty();
    var value = opt.value;
    console.log("countys value: " + value);
    $.ajax({
        type: "GET",
        url: "/citysinfo?belongs=" + value,
        contentType: 'application/json;charset=utf-8', //设置请求头信息
        dataType: "json",
        success: function (result, status) {
            console.log("status: " + status)
            var infos = result["regioninfo"];
            if(infos == null){
                console.log("区县信息为空")
                return
            }
            for(let i = 0; i < infos.length; i++){
                var str = '<option value="' + infos[i]["region_code"] + '">' + infos[i]["region_fullname"] +'</option>'
                $("#countys").append(str);
            }
        },
        error: function (result) {
            console.log(result.msg);
        }
    });
}

function initCity(){
    var el = document.getElementById("provinces");
    var code = el.options[el.options.selectedIndex].value;
    console.log("省份：" + code)
    $("#citys").empty();
    var strSpace = '<option value="space"></option>';
    $("#citys").append(strSpace);

    $.ajax({
        type: "GET",
        url: "/citysinfo?belongs=" + code,
        contentType: 'application/json;charset=utf-8', //设置请求头信息
        dataType: "json",
        success: function (result, status) {
            console.log("status: " + status)
            var infos = result["regioninfo"];
            for(let i = 0; i < infos.length; i++){
                var str = '<option value="' + infos[i]["region_code"] + '">' + infos[i]["region_fullname"] +'</option>'
                $("#citys").append(str);
            }
            initCounty()
        },
        error: function (result) {
            console.log(result.msg);
        }
    });
}


function initCounty(){
    var el = document.getElementById("citys");
    var code = el.options[el.options.selectedIndex].value
    // var code = document.getElementById("citys").options[0].value;
    $("#countys").empty();

    $.ajax({
        type: "GET",
        url: "/citysinfo?belongs=" + code,
        contentType: 'application/json;charset=utf-8', //设置请求头信息
        dataType: "json",
        success: function (result, status) {
            console.log("status: " + status)
            var infos = result["regioninfo"];
            if(infos == null){
                console.log("区县信息为空")
                return
            }
            for(let i = 0; i < infos.length; i++){
                var str = '<option value="' + infos[i]["region_code"] + '">' + infos[i]["region_fullname"] +'</option>'
                $("#countys").append(str);
            }
        },
        error: function (result) {
            console.log(result.msg);
        }
    });
}

function add_validate(){
    var el = document.getElementById("provinces");
    var value = el.options[el.options.selectedIndex].value;
    if(value == "space"){
        alert("省份不能为空");
        return false;
    }

    el = document.getElementById("citys");
    value = el.options[el.options.selectedIndex].value;
    if(value == "space"){
        alert("城市(直辖市-区)不能为空");
        return false;
    }

    //检验
    if(document.forms["modify-form"].storename.value == ""){
        alert("名称不能为空");
        return false;
    }
    if(document.forms["modify-form"].phone.value == ""){
        alert("电话不能为空");
        return false;
    }
    if(document.forms["modify-form"].province.value == ""){
        alert("省份不能为空");
        return false;
    }
    if(document.forms["modify-form"].city.value == ""){
        alert("城市(直辖市-区)不能为空");
        return false;
    }
    if(document.forms["modify-form"].county.value == ""){
        console.log("区县为空");
    }
    if(document.forms["modify-form"].address.value == ""){
        alert("地址不能为空");
        return false;
    }
    if(document.forms["modify-form"].tag.value == ""){
        alert("类别不能为空");
        return false;
    }

    return true
}

