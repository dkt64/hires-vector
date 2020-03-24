<template>
  <v-container>
    <v-card elevation="5" class="mx-auto" outlined>
      <v-list-item three-line>
        <v-list-item-content>
          <div class="overline mb-4">VIEW</div>
          <v-list-item-title class="headline mb-1">Output window</v-list-item-title>
          <v-list-item-subtitle>Here we can see input render</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-row>
        <v-col>
          <v-list-item-avatar class="mx-2" tile id="container" style="width: 640px; height: 400px"></v-list-item-avatar>
        </v-col>
        <v-col>
          <v-row>
            <v-switch v-model="switch_controls" label="Take over angle control" class="mb-5"></v-switch>
          </v-row>
          <v-row>
            <v-slider
              class="headline mr-5 mb-5"
              :thumb-size="24"
              thumb-label="always"
              v-model="slider_x"
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
              min="0"
              max="360"
              label="Angle Z"
            ></v-slider>
          </v-row>
        </v-col>
      </v-row>
    </v-card>
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
    switch_controls: false
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
      container.appendChild(this.renderer.domElement);
    },
    animate: function() {
      requestAnimationFrame(this.animate);
      this.izo.rotation.x += 0.01;
      this.izo.rotation.y += 0.01;
      this.renderer.render(this.scene, this.camera);
    }
  },
  mounted() {
    console.log("Mounted()");
    this.init();
    this.animate();
  }
};
</script>
