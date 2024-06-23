<?php

 

define("MB", 1048576);

function filterRequest($requestname)
{
    return  htmlspecialchars(strip_tags($_POST[$requestname]));
}
function sha1Request($requestname)
{
    return  sha1($_POST[$requestname]);
}
function getAllData($table, $where = null, $values = null, $json = true)
{
    global $con;
    $data = array();
    if ($where == null) {
        $stmt = $con->prepare("SELECT  * FROM  $table ");
    } else {
        $stmt = $con->prepare("SELECT  * FROM  $table WHERE   $where ");
    }

    $stmt->execute($values);
    $data = $stmt->fetchAll(PDO::FETCH_ASSOC);
    $count  = $stmt->rowCount();


    if ($json == true) {
        if ($count > 0) {
            echo json_encode(array("status" => "success", "data" => $data));
        } else {
            echo json_encode(array("status" => "failure"));
        }

        return $count;
    } else {
        if ($count > 0) {
            return $data;
        } else {
            return json_encode(array("status" => "failure"));
        }
    }
}

function getData($table, $where = null, $values = null, $json = true)
{
    global $con;
    $data = array();
    $stmt = $con->prepare("SELECT  * FROM $table WHERE   $where ");
    $stmt->execute($values);
    $data = $stmt->fetch(PDO::FETCH_ASSOC);
    $count  = $stmt->rowCount();
    // if ($count > 0) {
    //     echo json_encode(array("status" => "success", "data" => $data));
    // } else {
    //     echo json_encode(array("status" => "failure"));
    // }
    // return $count;

    if ($json == true) {
        if ($count > 0) {
            echo json_encode(array("status" => "success", "data" => $data));
        } else {
            echo json_encode(array("status" => "failure"));
        }

        return $count;
    } else {
        if ($count > 0) {
            return $data;
        } else {
            return json_encode(array("status" => "failure"));
        }
    }
}



function insertData($table, $data, $json = true)
{
    global $con;
    foreach ($data as $field => $v)
        $ins[] = ':' . $field;
    $ins = implode(',', $ins);
    $fields = implode(',', array_keys($data));
    $sql = "INSERT INTO $table ($fields) VALUES ($ins)";

    $stmt = $con->prepare($sql);
    foreach ($data as $f => $v) {
        $stmt->bindValue(':' . $f, $v);
    }
    $stmt->execute();
    $count = $stmt->rowCount();
    if ($json == true) {
        if ($count > 0) {
            echo json_encode(array("status" => "success"));
        } else {
            echo json_encode(array("status" => "failure"));
        }
    }
    return $count;
}


function updateData($table, $data, $where, $json = true)
{
    global $con;
    $cols = array();
    $vals = array();

    foreach ($data as $key => $val) {
        $vals[] = "$val";
        $cols[] = "`$key` =  ? ";
    }
    $sql = "UPDATE $table SET " . implode(', ', $cols) . " WHERE $where";

    $stmt = $con->prepare($sql);
    $stmt->execute($vals);
    $count = $stmt->rowCount();
    if ($json == true) {
        if ($count > 0) {
            echo json_encode(array("status" => "success"));
        } else {
            echo json_encode(array("status" => "failure", "data" => $data, "stmt" => $stmt));
        }
    }
    return $count;
}

function deleteData($table, $where, $json = true)
{
    global $con;
    $stmt = $con->prepare("DELETE FROM $table WHERE $where");
    $stmt->execute();
    $count = $stmt->rowCount();
    if ($json == true) {
        if ($count > 0) {
            echo json_encode(array("status" => "success"));
        } else {
            echo json_encode(array("status" => "failure"));
        }
    }
    return $count;
}

function imageUpload($dir, $imageRequest)
{
    if (isset($_FILES[$imageRequest])) {
        global $msgError;
        $imagename  = rand(1000, 100000) . $_FILES[$imageRequest]['name'];
        $imagetmp   = $_FILES[$imageRequest]['tmp_name'];
        $imagesize  = $_FILES[$imageRequest]['size'];
        $allowExt   = array("jpg", "png", "gif", "mp3", "pdf", "svg");
        $strToArray = explode(".", $imagename);
        $ext        = end($strToArray);
        $ext        = strtolower($ext);

        if (!empty($imagename) && !in_array($ext, $allowExt)) {
            $msgError = "EXT";
        }
        if ($imagesize > 2 * MB) {
            $msgError = "size";
        }
        if (empty($msgError)) {
            // move_uploaded_file($imagetmp,  "../upload/" . $imagename);
            move_uploaded_file($imagetmp, $dir . "/" . $imagename);

            return $imagename;
        } else { print($_FILES[$imageRequest]);
            return "fail";
        }
    } else {
        return "empty";
    }
}



function deleteFile($dir, $imagename)
{
    if (file_exists($dir . "/" . $imagename)) {
        unlink($dir . "/" . $imagename);
    }
}

function checkAuthenticate()
{
    if (isset($_SERVER['PHP_AUTH_USER'])  && isset($_SERVER['PHP_AUTH_PW'])) {
        if ($_SERVER['PHP_AUTH_USER'] != "wael" ||  $_SERVER['PHP_AUTH_PW'] != "wael12345") {
            header('WWW-Authenticate: Basic realm="My Realm"');
            header('HTTP/1.0 401 Unauthorized');
            echo 'Page Not Found';
            exit;
        }
    } else {
        exit;
    }

    // End 
}


function   printFailure($message = "none")
{
    echo     json_encode(array("status" => "failure", "message" => $message));
}
function   printsuccess($message = "none")
{
    echo     json_encode(array("status" => "success", "message" => $message));
}
function result($count)
{
    if ($count > 0) {
        printsuccess();
    } else {
        printFailure();
    }
}
function sendEmail($to, $title, $body)
{
    $header = "From: ahmedmady@amadytech.com ";
    // $header = "From: ahmedmady@amadytech.com " . "\n" . "CC:ahmed.eid.ac.1.edu@gmail.com";
    mail($to, $title, $body, $header);
    // echo json_encode(["status" => "Success"]);
}
/////////////////////////////////////////////////////////////////////////
function sendGCM($title, $message, $topic, $pageid, $pagename)
{


    $url = 'https://fcm.googleapis.com/fcm/send';

    $fields = array(
        "to" => '/topics/' . $topic,
        'priority' => 'high',
        'content_available' => true,

        'notification' => array(
            "body" =>  $message,
            "title" =>  $title,
            "click_action" => "FLUTTER_NOTIFICATION_CLICK",
            "sound" => "default"

        ),
        'data' => array(
            "pageid" => $pageid,
            "pagename" => $pagename
        )

    );


    $fields = json_encode($fields);
    $headers = array(
        'Authorization: key=' . 
        "AAAAF-G7T6A:APA91bECGk1HE7RnPKQq50_R8uryP08uZVPYK5lsaGlCwi2k5MHFDursMHtB-WnHMMk_qMlKTJSleprGRkbr_7qC2AX-UHxyKkj5_RsHPWGOGF--XiN8uBXC5AmKiz5eDxavE4NUJDf_",
        'Content-Type: application/json'
    );
            // "AAAAwDEe1lE:APA91bHjnYAADGH4fINZZZxB4IS7FnSqX6-
        // ig44k2bRN4jTU4dTNSYCP3DprTEsyZDDXjd0A1Dhkzv2HEm5u2bGpDNEP7LNo
        // -KqDQ6oIjHmbiFbtRPBbo8UghSAdlTkc3EjHDR7ozRvC",

    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $fields);

    $result = curl_exec($ch);
    return $result;
    curl_close($ch);
}


function insertNotify(
    $userid,
    $title,
    $body,
    $topic,
    $pageid,
    $pagename
) {
    global $con;
    $stmt = $con->prepare("INSERT INTO `notifications`( `notifications_userid`, `notifications_title`, `notifications_body`,`notifications_topic`)
    VALUES ( ?,?,?,?)");
    $stmt->execute(array($userid, $title, $body, $topic));
    // $stmt->execute(array($topic, $title, $body));//for show notif. for "users" topic also
    sendGCM($title, $body, $topic, $pageid, $pagename);
    $count = $stmt->rowCount();
    return $count;
}
