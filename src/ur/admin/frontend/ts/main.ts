import * as angular from "angular";

namespace UrAdmin {

    let app = angular.module('AdminApp', ['ngMaterial']);

    app.controller("DefaultCtrl", function($scope) {
        console.log("hello");
    });
}
