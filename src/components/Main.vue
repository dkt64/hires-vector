<template>
  <v-container>
    <v-row>
      <v-row>
        <v-col>
          <v-list-item-title class="headline mb-1">Input window</v-list-item-title>
          <div id="container" style="width: 320px; height: 200px"></div>
        </v-col>
        <v-col>
          <v-list-item-title class="headline mb-1">Hires window</v-list-item-title>
          <canvas id="canvasHires" width="320" height="200"></canvas>
        </v-col>
        <v-col>
          <v-row>
            <v-switch v-model="switch_controls" label="Animate" class="overline mb-5"></v-switch>
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
    pixels: null,
    glContext: null,
    hiresData: null,
    canvasHires: null,
    contextHires: null
  }),
  methods: {
    init: function() {
      let container = document.getElementById("container");

      this.camera = new THREE.PerspectiveCamera(
        60,
        container.clientWidth / container.clientHeight,
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
          0xa050a0 + (i * 4 + 1) + (i * 4 + 1) * 256 + (i * 4 + 1) * 256 * 256
        );
      }
      this.izo = new THREE.Mesh(izoGeo, izoMat);
      this.izo.name = "izo";
      console.log(this.izo.name);
      console.log(this.izo);
      this.scene.add(this.izo);

      this.renderer = new THREE.WebGLRenderer({ antialias: true });
      this.renderer.setSize(container.clientWidth, container.clientHeight);
      this.renderer.setScissor(0, 0, 320, 200);
      this.renderer.setScissorTest(true);
      container.appendChild(this.renderer.domElement);

      this.renderer.render(this.scene, this.camera);

      this.glContext = this.renderer.getContext();
      this.pixels = new Uint8Array(320 * 200 * 4);
      this.glContext.readPixels(
        0,
        0,
        320,
        200,
        this.glContext.RGBA,
        this.glContext.UNSIGNED_BYTE,
        this.pixels
      );
      console.log("pixels:");
      console.log(this.pixels); // Uint8Array

      this.canvasHires = document.getElementById("canvasHires");
      if (this.canvasHires.getContext) {
        this.contextHires = this.canvasHires.getContext("2d");
        console.log("Got 2d canvas contex:");
        console.log(this.contextHires);
      } else {
        console.log("No 2d canvas. Sorry...");
      }
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

      this.glContext.readPixels(
        0,
        0,
        320,
        200,
        this.glContext.RGBA,
        this.glContext.UNSIGNED_BYTE,
        this.pixels
      );

      this.hiresData = this.contextHires.getImageData(0, 0, 320, 200);

      var k = 0;
      for (var j = 200; j > 0; --j) {
        for (var i = 0; i < 320 * 4; i++) {
          this.hiresData.data[k++] = this.pixels[j * 320 * 4 + i];
        }
      }

      this.contextHires.putImageData(this.hiresData, 0, 0);
    }
  },
  mounted() {
    console.log("Mounted()");
    this.init();
    this.animate();
  }
};
</script>
