<template>
  <v-container>
    <v-row class="d-flex justify-space-around">
      <v-card tile outlined class="mx-auto pa-2">
        <v-toolbar color="indigo" class="headline" dark>Three.js</v-toolbar>
        <div v-on:click="clickOnGL" id="canvasGL" style="width: 320px; height: 200px"></div>
      </v-card>
      <v-card tile outlined class="mx-auto pa-2">
        <v-toolbar color="indigo" class="headline" dark>Hires</v-toolbar>
        <canvas id="canvasHires" width="320" height="200"></canvas>
      </v-card>
      <v-card tile outlined class="mx-auto pa-2">
        <v-toolbar color="indigo" class="headline" dark>Diff</v-toolbar>
        <canvas id="canvasDiff" width="320" height="200"></canvas>
      </v-card>
    </v-row>
    <v-divider class="mt-5"></v-divider>
    <v-row class="mt-8 d-flex justify-space-around">
      <v-slider
        class="ml-5"
        :thumb-size="24"
        thumb-label="always"
        v-model="slider_x"
        step="1"
        min="0"
        max="360"
        label="Angle X"
      ></v-slider>
      <v-slider
        class="ml-5"
        :thumb-size="24"
        thumb-label="always"
        v-model="slider_y"
        step="1"
        min="0"
        max="360"
        label="Angle Y"
      ></v-slider>
      <v-slider
        class="ml-5 mr-2"
        :thumb-size="24"
        thumb-label="always"
        v-model="slider_z"
        step="1"
        min="0"
        max="360"
        label="Angle Z"
      ></v-slider>
    </v-row>
    <v-row class="d-flex justify-start">
      <v-switch v-model="switch_controls" label="Animate"></v-switch>
    </v-row>
    <v-row>
      <v-chip :text-color="backgroundColorText" :color="backgroundColorBack">Background color [mouse click]</v-chip>
    </v-row>
  </v-container>
</template>

<script>
import * as THREE from "three";

var pixelsGL = new Uint8Array(320 * 200 * 4);
var changeBackColor = false;
var backgroundColor = [0, 0, 0, 1];

export default {
  name: "Main",

  data: () => ({
    renderer: null,
    scene: null,
    camera: null,
    izo: null,
    slider_x: 0,
    slider_y: 0,
    slider_z: 0,
    switch_controls: false,
    backgroundColorBack: "#000000ff",
    backgroundColorText: "#ffffffff"
  }),
  methods: {
    colorString: function(r, g, b, a) {
      return (
        "rgba(" +
        r.toString() +
        "," +
        g.toString() +
        "," +
        b.toString() +
        "," +
        (255 / a).toString() +
        ")"
      );
    },
    clickOnGL: function(event) {
      var canvasGL = document.getElementById("canvasGL");

      console.log("mouse clicked ", event);
      console.log("this.canvasGL ", canvasGL);
      var x = event.layerX - canvasGL.offsetLeft;
      var y = event.layerY - canvasGL.offsetTop;
      console.log("mouse clicked at ", x, " ", y);

      var offset = x * 4 + y * 320 * 4;
      console.log("offset: " + offset);

      backgroundColor[0] = pixelsGL[offset + 0];
      backgroundColor[1] = pixelsGL[offset + 1];
      backgroundColor[2] = pixelsGL[offset + 2];
      backgroundColor[3] = pixelsGL[offset + 3];

      console.log(
        "New background color: " +
          this.colorString(
            backgroundColor[0],
            backgroundColor[1],
            backgroundColor[2],
            backgroundColor[3]
          )
      );

      changeBackColor = true;
    },
    init: function() {
      var canvasGL = document.getElementById("canvasGL");
      this.camera = new THREE.PerspectiveCamera(
        60,
        canvasGL.clientWidth / canvasGL.clientHeight,
        0.1,
        1000
      );
      this.camera.position.z = 200;

      this.scene = new THREE.Scene();

      const light = new THREE.AmbientLight(0xffffff);
      // light.position.set(0, 0, 1000);
      this.scene.add(light);

      var izoGeo = new THREE.IcosahedronGeometry(60);
      var izoMat = new THREE.MeshPhongMaterial({
        color: 0xffffff,
        vertexColors: THREE.FaceColors
      });
      for (var i = 0; i < izoGeo.faces.length; i++) {
        //izoGeo.faces[i].color.setHex((c64colors[i & 0x0f] * (i + 1)) & 0xffff00 + i);
        izoGeo.faces[i].color.setHex(
          0x202020 + (i * 8 + 1) + (i * 8 + 1) * 256 + (i * 8 + 1) * 256 * 256
        );
      }
      this.izo = new THREE.Mesh(izoGeo, izoMat);
      this.izo.name = "izo";
      console.log(this.izo.name);
      console.log(this.izo);
      this.scene.add(this.izo);

      this.renderer = new THREE.WebGLRenderer({ antialias: true });
      this.renderer.setSize(canvasGL.clientWidth, canvasGL.clientHeight);
      this.renderer.setScissor(0, 0, 320, 200);
      this.renderer.setScissorTest(true);
      canvasGL.appendChild(this.renderer.domElement);

      this.renderer.render(this.scene, this.camera);
    },
    animate: function() {
      requestAnimationFrame(this.animate);

      if (!this.switch_controls) {
        this.izo.rotation.x = (Math.PI / 180) * this.slider_x;
        this.izo.rotation.y = (Math.PI / 180) * this.slider_y;
        this.izo.rotation.z = (Math.PI / 180) * this.slider_z;
      } else {
        this.izo.rotation.x += 0.04;
        this.izo.rotation.y += 0.02;
        this.izo.rotation.z -= 0.02;

        this.slider_x = Math.abs(((this.izo.rotation.x / Math.PI) * 180) % 360);
        this.slider_y = Math.abs(((this.izo.rotation.y / Math.PI) * 180) % 360);
        this.slider_z = Math.abs(((this.izo.rotation.z / Math.PI) * 180) % 360);
      }
      this.renderer.render(this.scene, this.camera);

      var contextGL = this.renderer.getContext();

      contextGL.readPixels(
        0,
        0,
        320,
        200,
        contextGL.RGBA,
        contextGL.UNSIGNED_BYTE,
        pixelsGL
      );

      var canvasHires = document.getElementById("canvasHires");
      var contextHires = canvasHires.getContext("2d");

      var canvasDiff = document.getElementById("canvasDiff");
      var contextDiff = canvasDiff.getContext("2d");

      var pixelsHires = contextHires.getImageData(0, 0, 320, 200);
      var k = 0;
      for (var j = 200; j > 0; --j) {
        for (var i = 0; i < 320 * 4; i++) {
          pixelsHires.data[k++] = pixelsGL[j * 320 * 4 + i];
        }
      }

      // Hires
      //

      // Zapis do output
      contextHires.putImageData(pixelsHires, 0, 0);
      contextDiff.putImageData(pixelsHires, 0, 0);

      // Zmiana background color na podstawie kliknięcia myszką
      if (changeBackColor) {
        this.backgroundColorBack = this.colorString(
          backgroundColor[0],
          backgroundColor[1],
          backgroundColor[2],
          backgroundColor[3]
        );
        if (
          (backgroundColor[0] + backgroundColor[0] + backgroundColor[0]) / 3 <
          0x80
        ) {
          this.backgroundColorText = this.colorString(0xff, 0xff, 0xff, 0xff);
        } else {
          this.backgroundColorText = this.colorString(0, 0, 0, 0xff);
        }
      }
    }
  },
  mounted() {
    console.log("Mounted()");
    this.init();
    this.animate();
  }
};
</script>
