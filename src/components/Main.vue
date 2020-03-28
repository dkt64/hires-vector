<template>
  <v-container>
    <v-row>
      <v-col>
        <v-list-item-title class="headline">Input window</v-list-item-title>
        <div v-on:click="clickOnContainer" id="canvasGL" style="width: 320px; height: 200px"></div>
        <v-chip :text-color="backgroundColorText" :color="backgroundColorBack" class="mt-1">Background color</v-chip>
      </v-col>
      <v-col>
        <v-list-item-title class="headline">Hires window</v-list-item-title>
        <canvas id="canvasHires" width="320" height="200"></canvas>
      </v-col>
      <v-col>
        <v-row>
          <v-switch v-model="switch_controls" label="Animate" class="headline"></v-switch>
        </v-row>
        <v-row>
          <v-slider
            class="headline mr-5 mb-5"
            :thumb-size="24"
            thumb-label="always"
            v-model="slider_x"
            step="1"
            min="0"
            max="360"
            label="Angle X"
          ></v-slider>
        </v-row>
        <v-row>
          <v-slider
            class="headline mr-5 mb-5"
            :thumb-size="24"
            thumb-label="always"
            v-model="slider_y"
            step="1"
            min="0"
            max="360"
            label="Angle Y"
          ></v-slider>
        </v-row>
        <v-row>
          <v-slider
            class="headline mr-5"
            :thumb-size="24"
            thumb-label="always"
            v-model="slider_z"
            step="1"
            min="0"
            max="360"
            label="Angle Z"
          ></v-slider>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import * as THREE from "three";

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
    pixelsGL: null,
    pixelsHires: null,
    canvasGL: null,
    canvasHires: null,
    contextGL: null,
    contextHires: null,
    backgroundColor: [0, 0, 0, 1],
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
    clickOnContainer: function(event) {
      console.log("mouse clicked ", event);
      console.log("this.canvasGL ", this.canvasGL);
      var x = event.layerX - this.canvasGL.offsetLeft;
      var y = event.layerY - this.canvasGL.offsetTop;
      console.log("mouse clicked at ", x, " ", y);

      var offset = x * 4 + y * 320 * 4;

      this.backgroundColor[0] = this.pixelsGL[offset + 0];
      this.backgroundColor[1] = this.pixelsGL[offset + 1];
      this.backgroundColor[2] = this.pixelsGL[offset + 2];
      this.backgroundColor[3] = this.pixelsGL[offset + 3];

      this.backgroundColorBack = this.colorString(
        this.backgroundColor[0],
        this.backgroundColor[1],
        this.backgroundColor[2],
        this.backgroundColor[3]
      );

      this.backgroundColorText = this.colorString(
        this.backgroundColor[0] ^ 0x80,
        this.backgroundColor[1] ^ 0x80,
        this.backgroundColor[2] ^ 0x80,
        this.backgroundColor[3]
      );

      console.log(
        "New background color: " +
          this.colorString(
            this.backgroundColor[0],
            this.backgroundColor[1],
            this.backgroundColor[2],
            this.backgroundColor[3]
          )
      );
    },
    init: function() {
      this.canvasGL = document.getElementById("canvasGL");

      this.camera = new THREE.PerspectiveCamera(
        60,
        this.canvasGL.clientWidth / this.canvasGL.clientHeight,
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
      this.renderer.setSize(
        this.canvasGL.clientWidth,
        this.canvasGL.clientHeight
      );
      this.renderer.setScissor(0, 0, 320, 200);
      this.renderer.setScissorTest(true);
      this.canvasGL.appendChild(this.renderer.domElement);

      this.renderer.render(this.scene, this.camera);

      this.contextGL = this.renderer.getContext();
      this.pixelsGL = new Uint8Array(320 * 200 * 4);
      this.contextGL.readPixels(
        0,
        0,
        320,
        200,
        this.contextGL.RGBA,
        this.contextGL.UNSIGNED_BYTE,
        this.pixelsGL
      );
      console.log("pixelsGL:");
      console.log(this.pixelsGL); // Uint8Array

      this.canvasHires = document.getElementById("canvasHires");
      if (this.canvasHires.getContext) {
        this.contextHires = this.canvasHires.getContext("2d");
        console.log("Got 2d canvas contex:");
        console.log(this.contextHires);
      } else {
        console.log("No 2d canvas. Sorry...");
      }

      this.backgroundColor[0] = this.pixelsGL[0];
      this.backgroundColor[1] = this.pixelsGL[1];
      this.backgroundColor[2] = this.pixelsGL[2];
      this.backgroundColor[3] = this.pixelsGL[3];
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

      this.contextGL.readPixels(
        0,
        0,
        320,
        200,
        this.contextGL.RGBA,
        this.contextGL.UNSIGNED_BYTE,
        this.pixelsGL
      );

      this.pixelsHires = this.contextHires.getImageData(0, 0, 320, 200);

      var k = 0;
      for (var j = 200; j > 0; --j) {
        for (var i = 0; i < 320 * 4; i++) {
          this.pixelsHires.data[k++] = this.pixelsGL[j * 320 * 4 + i];
        }
      }

      this.contextHires.putImageData(this.pixelsHires, 0, 0);
    }
  },
  mounted() {
    console.log("Mounted()");
    this.init();
    this.animate();
  },
};
</script>
