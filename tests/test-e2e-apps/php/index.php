<?php
use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Slim\Factory\AppFactory;

require __DIR__ . '/vendor/autoload.php';

$app = AppFactory::create();

$app->get('/', function (Request $request, Response $response) {
    $all_envs = getenv();
    $all_values_string = implode(", ", $all_envs);
    $response->getBody()->write(all_values_string);
    // $result = random_int(1,6);
    // $response->getBody()->write(strval($result));
    return $response;
});

$app->run();
