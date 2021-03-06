var scene3d = document.getElementById("scene3d");
var scene3dW = scene3d.offsetWidth;
var scene3dH = scene3d.offsetHeight;
     
var scene = new THREE.Scene();
var camera = new THREE.PerspectiveCamera( 75, scene3dW / scene3dW, 0.1, 1000 );

var playerDirection = 0; // angles 0 - 2pi
var dVector;
var angularSpeed = 0.01;
var angularSpeed = 0.01;
var playerSpeed = 0.075;
var playerBackwardsSpeed = playerSpeed * 0.4;

// 2 axes for roll, pitch, yaw
// var myCameraRollAxis = new THREE.Vector3( 0, 0, 1 );
// var myCameraPitchAxis = new THREE.Vector3( 0, 0, 0 );
// var myCameraYawAxis = new THREE.Vector3( 0, 0, 0 );


var renderer = new THREE.WebGLRenderer();
// renderer.setSize( window.innerWidth, window.innerHeight );
renderer.setSize( scene3dW, scene3dH );

// document.body.appendChild( renderer.domElement );
scene3d.appendChild(renderer.domElement);

var geometry = new THREE.BoxGeometry( 1, 0.1, 1 );

camera.position.x = 0;
camera.position.y = 12;
camera.position.z = 1;

// dVector = new THREE.Vector3( 0, 0, 1 ) ;
// camera.lookAt( dVector );

var controls = new THREE.OrbitControls(camera, renderer.domElement);

var playerIsRotatingLeft = 0;
var playerIsRotatingRight = 0;
var playerIsMovingForward = 0;
var playerIsMovingBackwards = 0;
    
function key_down(event){
var W = 87;
var S = 83;
var A = 65;
var D = 68;
var minus = 189;
var plus = 187;

var k = event.keyCode;
console.log(k);
if (k == A){ // rotate left
    playerIsRotatingLeft = 1;
}
if (k == D){ // rotate right
    playerIsRotatingRight = 1;
}
if (k == W){ // go forward
    playerIsMovingForward = 1;
}
if (k == S){ // go back 
    playerIsMovingBackwards = 1;
}
}

document.addEventListener("keydown", key_down, false);
document.addEventListener("keyup", key_up, false);

// Create field of flat tiles.
var materialField1 = new THREE.MeshBasicMaterial( { color: 0xbbbb00 } );

var materialField2 = new THREE.MeshBasicMaterial( { color: 0xbb0000 } );

var numX = 5;
var numY = 5;
var fieldElements = new Array(numX);
for (var x = 0; x < numX; x++) { 
fieldElements[x] = new Array(numY);
for (var y = 0; y < numY; y++) {
    if ((x+y) % 2 == 0)
        fieldElements[x][y] = new THREE.Mesh( geometry, materialField1 );
    else
        fieldElements[x][y] = new THREE.Mesh( geometry, materialField2 );
    fieldElements[x][y].position.set( x * 2 - 0, 0, y * 2 - 0 );
    scene.add( fieldElements[x][y] );
}
}

// The X axis is red. The Y axis is green. The Z axis is blue.
var axesHelper = new THREE.AxesHelper( 5 );
scene.add( axesHelper );

var counter = 0;

function setPlayerDirection(){

return;
// if (counter == 0)

var delta_x = playerSpeed * Math.cos(playerDirection);
var delta_z = playerSpeed * Math.sin(playerDirection);

console.log("camera.position.x: " + camera.position.x);
console.log("camera.position.y: " + camera.position.y);
console.log("camera.position.z: " + camera.position.z);

var new_dx = camera.position.x + delta_x;
var new_dz = camera.position.z + delta_z;
dVector.x = new_dx;
dVector.z = new_dz;

console.log(dVector);

camera.lookAt( dVector ); 
}

function updatePlayer(){
if (playerIsRotatingLeft){ // rotate left
    playerDirection -= angularSpeed;
}
if (playerIsRotatingRight){ // rotate right
    playerDirection += angularSpeed;
}
if (playerIsRotatingRight || playerIsRotatingLeft){
    setPlayerDirection();
    return;
}
if (playerIsMovingForward){ // go forward
    moveForward(playerSpeed);
    return;
}
if (playerIsMovingBackwards){ // go backwards
    moveForward(-playerBackwardsSpeed);
    return;
}

}

function spinTiles() {
for (var x = 0; x < numX; x++) { 
    for (var y = 0; y < numY; y++) {
        fieldElements[x][y].rotation.x += 0.01;
        fieldElements[x][y].rotation.y += 0.01;
    }
}
}

function animate() {
requestAnimationFrame( animate );

spinTiles();

updatePlayer();

renderer.render( scene, camera );
}
animate();

function key_up(event){
playerIsMovingForward = 0;
playerIsMovingBackwards = 0;
playerIsRotatingLeft = 0;
playerIsRotatingRight = 0;
playerGoesUp = 0;
playerGoesDown = 0;
}

function key_down(event){
var W = 87;
var S = 83;
var A = 65;
var D = 68;
var minus = 189;
var plus = 187;

var k = event.keyCode;
console.log(k);
if(k == A){ // rotate left
    playerIsRotatingLeft = 1;
}
if(k == D){ // rotate right
    playerIsRotatingRight = 1;
}
if(k == W){ // go forward
    playerIsMovingForward = 1;
}
if(k == S){ // go back 
    playerIsMovingBackwards = 1;
}
}

function moveForward(speed){
var delta_x = speed * Math.cos(playerDirection);
var delta_z = speed * Math.sin(playerDirection);
var new_x = camera.position.x + delta_x;
var new_z = camera.position.z + delta_z;
camera.position.x = new_x;
camera.position.z = new_z;

var new_dx = dVector.x + delta_x;
var new_dz = dVector.z + delta_z;
dVector.x = new_dx;
dVector.z = new_dz;
camera.lookAt( dVector );    

}
