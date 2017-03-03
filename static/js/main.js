var server_manager = angular.module('server_manager', []).controller("mainController");


function mainController($scope, $http) {
    // when landing on the page, get all torrents and show them? Naw....
    // TODO: SHould be a get with the url as parameter

    generic_fail_funciton = function(data, status, headers, config) {
        $scope.status = status + ' ' + headers;
    }

    send_server_command = function(data, success_funciton, fail_function) {
        fail_function = fail_function || generic_fail_funciton;
        $http({
            url: '/api/cmd',
            method: "POST",
            data: JSON.stringify(data),
            headers: {'Content-Type': 'application/json'}
        }).success(success_funciton)
          .error(fail_function);
    }

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
        data = {"action" : "list_dir", "args":{ "dir" : "/Users/erikparreira/Dropbox/Programming/home_server_manager/test_data"}};
        send_server_command(data, function(data, status, headers, config) {
            console.log("DONE!");
            console.log(data);
            $scope.files = data.Args;
        });
    }

    $scope.files = ["hello"];

}

server_manager.controller("mainController", mainController);