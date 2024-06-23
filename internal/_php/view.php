<?php

include "../connect.php";
$userid=filterRequest("userid");
getAllData("address","address_userid=?",array($userid));