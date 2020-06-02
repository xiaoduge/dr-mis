
function queryStore(){
    console.log("查询门店信息")
    var name = document.forms["modify-form"].storename.value;
    var phone = document.forms["modify-form"].phone.value;

    $.ajax({
        type: "GET",
        url: "/mis/data/query-store?storename=" + name + "&storephone=" + phone,
        contentType: 'application/json;charset=utf-8', //设置请求头信息
        dataType: "json",
        success: function (result, status) {
            if (result["store_address"] == ""){
                document.forms["modify-form"].storename.value = "没有该门店的信息";
                return
            }
            console.log("query status: " + status)
            document.forms["modify-form"].storeid.value = result["store_id"]; //获取一个ID
            document.forms["modify-form"].storename.value = result["store_name"];
            document.forms["modify-form"].phone.value = result["store_phone"];
            document.forms["modify-form"].address.value = result["store_address"];
            document.forms["modify-form"].mark.value = result["store_tag"];
            // 设置省份
            $("#provinces").empty();
            var str = '<option value="' + result["store_province_code"] + '">' + result["store_province"]+'</option>'
            $("#provinces").append(str);
            //
            $("#citys").empty();
            var str = '<option value="' + result["store_city_code"] + '">' + result["store_city"]+'</option>'
            $("#citys").append(str);
            //
            $("#countys").empty();
            var str = '<option value="' + result["store_county_code"] + '">' + result["store_county"]+'</option>'
            $("#countys").append(str);

        },
        error: function (result) {
            console.log(result.msg);
        }
    });
}

function submitStore(){
    console.log("更新门店信息")
    if(add_validate()){
        document.forms["modify-form"].submit();
    }
}

function deleteStore(){
    console.log("删除门店信息");
    var id = document.forms["modify-form"].storeid.value;
    console.log("id: " + id);
    $.ajax({
        type: "GET",
        url: "/mis/data/delete-storeinfo?storeid=" + id,
        contentType: 'application/json;charset=utf-8', //设置请求头信息
        dataType: "json",
        success: function (result, status) {
            console.log("delete status: " + status);
            document.forms["modify-form"].storeid.value = result["status"];
        },
        error: function (result) {
            console.log("delete store error: " + result.msg);
        }
    });
}

function editCounty(){
    initCounty();
}

function editCity(){
    initCity();
}

function editProvince(){
    initProvince();
}