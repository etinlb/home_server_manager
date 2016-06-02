var server_manager = angular.module('server_manager', []).controller("mainController");


function mainController($scope, $http) {
    // when landing on the page, get all torrents and show them? Naw....
    // TODO: SHould be a get with the url as parameter
    $http({
            url: '/api/cmd',
            method: "POST",
            data: JSON.stringify( {"action" : "list_dir", "args":{ "dir" : "/Users/erikparreira/Dropbox/Programming/home_server_manager/test_data"}}),
            headers: {'Content-Type': 'application/json'}
        }).success(function (data, status, headers, config) {
            console.log("DONE!");
            console.log(data);
        }).error(function (data, status, headers, config) {
            $scope.status = status + ' ' + headers;
        });

     $scope.fix_permissions = function() {
        $http({
            url: '/api/cmd',
            method: "POST",
            data: JSON.stringify({args:{ "dir" : "fix_all_permissions"}}),
            headers: {'Content-Type': 'application/json'}
        }).success(function (data, status, headers, config) {
            console.log("DONE!")
        }).error(function (data, status, headers, config) {
            $scope.status = status + ' ' + headers;
        });
    }

    $scope.get_dir = function() {
        $http({
            url: '/api/cmd',
            method: "POST",
            data: JSON.stringify({action:"fix_all_permissions"}),
            headers: {'Content-Type': 'application/json'}
        }).success(function (data, status, headers, config) {
            console.log("DONE!")
        }).error(function (data, status, headers, config) {
            $scope.status = status + ' ' + headers;
        });
    }

}

server_manager.controller("mainController", mainController);