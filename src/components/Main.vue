<template>
  <v-container>
    <v-row>
      <v-col>
        <v-list-item-title class="headline">Input window</v-list-item-title>
        <div v-on:click="clickOnContainer" id="canvasGL" style="width: 320px; height: 200px"></div>
        <v-chip
          :text-color="backgroundColorText"
          :color="backgroundColorBack"
          class="mt-1"
        >Background color</v-chip>
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

      var canvasGL = document.getElementById("canvasGL");

      var contextGL = this.renderer.getContext();

      var pixelsGL = new Uint8Array(320 * 200 * 4);

      contextGL.readPixels(
        0,
        0,
        320,
        200,
        contextGL.RGBA,
        contextGL.UNSIGNED_BYTE,
        pixelsGL
      );

      console.log("mouse clicked ", event);
      console.log("this.canvasGL ", canvasGL);
      var x = event.layerX - canvasGL.offsetLeft;
      var y = event.layerY - canvasGL.offsetTop;
      console.log("mouse clicked at ", x, " ", y);

      var offset = x * 4 + y * 320 * 4;

      this.backgroundColor[0] = pixelsGL[offset + 0];
      this.backgroundColor[1] = pixelsGL[offset + 1];
      this.backgroundColor[2] = pixelsGL[offset + 2];
      this.backgroundColor[3] = pixelsGL[offset + 3];

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
      this.renderer.setSize(
        canvasGL.clientWidth,
        canvasGL.clientHeight
      );
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

      var pixelsGL = new Uint8Array(320 * 200 * 4);

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

      var pixelsHires = contextHires.getImageData(0, 0, 320, 200);
      var k = 0;
      for (var j = 200; j > 0; --j) {
        for (var i = 0; i < 320 * 4; i++) {
          pixelsHires.data[k++] = pixelsGL[j * 320 * 4 + i];
        }
      }

      contextHires.putImageData(pixelsHires, 0, 0);
    }
  },
  mounted() {
    console.log("Mounted()");
    this.init();
    this.animate();
  }
};
</script>
