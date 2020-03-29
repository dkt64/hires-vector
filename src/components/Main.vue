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
      <v-chip
        :text-color="backgroundColorText"
        :color="backgroundColorBack"
      >Background color [mouse click]</v-chip>
    </v-row>
  </v-container>
</template>

<script>
// ================================================================================================
// Importy oraz globalne zmienne
// ================================================================================================
import * as THREE from "three";

var pixelsGL = new Uint8Array(320 * 200 * 4);
var changeBackColor = false;
var backgroundColorRGBA = [0, 0, 0, 1];

// Przygotowanie buforów na pixele
// ------------------------------------------------------------------------------------------
// var pixelsHires = contextHires.getImageData(0, 0, 320, 200);
var pixelsGLInverted = new ImageData(320, 200);
var pixelsHires = new ImageData(320, 200);
var pixelsDiff = new ImageData(320, 200);

export default {
  name: "Main",

  // ================================================================================================
  // Zmienne VUE
  // ================================================================================================
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

  // ================================================================================================
  // Moje funkcje
  // ================================================================================================
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

    // --------------------------------------------------------------------------------------------
    // clickOnGL - kliknięcie na context GL
    // --------------------------------------------------------------------------------------------
    clickOnGL: function(event) {
      var canvasGL = document.getElementById("canvasGL");

      console.log("mouse clicked ", event);
      console.log("this.canvasGL ", canvasGL);
      var x = event.layerX - canvasGL.offsetLeft;
      var y = event.layerY - canvasGL.offsetTop;
      y = 200 - y;
      console.log("mouse clicked at ", x, " ", y);

      var offset = x * 4 + y * 320 * 4;
      console.log("offset: " + offset);

      backgroundColorRGBA[0] = pixelsGL[offset + 0];
      backgroundColorRGBA[1] = pixelsGL[offset + 1];
      backgroundColorRGBA[2] = pixelsGL[offset + 2];
      backgroundColorRGBA[3] = pixelsGL[offset + 3];

      console.log(
        "New background color: " +
          this.colorString(
            backgroundColorRGBA[0],
            backgroundColorRGBA[1],
            backgroundColorRGBA[2],
            backgroundColorRGBA[3]
          )
      );

      changeBackColor = true;
    },
    // --------------------------------------------------------------------------------------------
    // init - wywoływana w Mounted()
    // --------------------------------------------------------------------------------------------
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
      for (let i = 0; i < izoGeo.faces.length; i++) {
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

      this.renderer = new THREE.WebGLRenderer({ antialias: false });
      this.renderer.setSize(canvasGL.clientWidth, canvasGL.clientHeight);
      this.renderer.setScissor(0, 0, 320, 200);
      this.renderer.setScissorTest(true);
      canvasGL.appendChild(this.renderer.domElement);

      this.renderer.render(this.scene, this.camera);
    },
    // --------------------------------------------------------------------------------------------
    // animate - główna pętla renderingu
    // --------------------------------------------------------------------------------------------
    animate: function() {
      // Obsługa Three.js
      // ------------------------------------------------------------------------------------------
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

      // Pobranie kontekstu i odczyt pixeli
      // ------------------------------------------------------------------------------------------
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

      // Obracamy do góry nogami (pixele GL liczone są w Y od dołu)
      // ------------------------------------------------------------------------------------------
      var k = 0;
      for (let j = 199; j >= 0; j--) {
        for (let i = 0; i < 320 * 4; i++) {
          pixelsGLInverted.data[k++] = pixelsGL[j * 320 * 4 + i];
        }
      }

      // Konwersja na Hires
      // ------------------------------------------------------------------------------------------
      if (changeBackColor) {

        console.log("pixelsGLInverted:")
        console.log(pixelsGLInverted)

        // do celów testowych
        // Wypełniamy kanał ALFA
        // ------------------------------------------------------------------------------------------
        for (let j = 0; j < 200; j++) {
          for (let i = 0; i < 320 * 4; i++) {
            pixelsHires.data[j * 320 * 4 + i * 4 + 0] = backgroundColorRGBA[0];
            pixelsHires.data[j * 320 * 4 + i * 4 + 1] = backgroundColorRGBA[1];
            pixelsHires.data[j * 320 * 4 + i * 4 + 2] = backgroundColorRGBA[2];
            pixelsHires.data[j * 320 * 4 + i * 4 + 3] = 0xff;
            pixelsDiff.data[j * 320 * 4 + i * 4 + 0] = backgroundColorRGBA[0];
            pixelsDiff.data[j * 320 * 4 + i * 4 + 1] = backgroundColorRGBA[1];
            pixelsDiff.data[j * 320 * 4 + i * 4 + 2] = backgroundColorRGBA[2];
            pixelsDiff.data[j * 320 * 4 + i * 4 + 3] = 0xff;
          }
        }

        // Szukamy obszarów 8x8
        //
        for (let by = 100; by < 108; by += 8) {
          for (let bx = 320/4; bx < 320*3/4; bx += 8) {
            var foundBackColor = false;
            var newColor1 = backgroundColorRGBA[0];
            var newColor2 = backgroundColorRGBA[0];

            // Pierwsza pętla szuka kolorów dla Hires
            // Sprawdzamy tylko pierwszy bajt z RGBA
            //
            for (let sy = 0; sy < 8; sy++) {
              for (let sx = 0; sx < 8; sx++) {
                // Odczyt pixela
                //
                let pixelColor =
                  pixelsGLInverted.data[
                    (by + sy) * 320 * 4 + bx * 8 * 4 + sx * 4
                  ];

                // Sprawdzamy czy background
                // Jeżeli tak to zaznaczamy że znaleźliśmy
                // Jeżeli nie to zapisujemy nowy kolor
                if (pixelColor == backgroundColorRGBA[0]) {
                  foundBackColor = true;
                  // console.log("znalazłem tło");
                } else if (pixelColor != newColor1) {
                  newColor1 = pixelColor;
                  // console.log("znalazłem color 1");
                } else if (pixelColor != newColor2) {
                  // console.log("znalazłem color 2");
                  newColor2 = pixelColor;
                }
              }
            }

            // Druga pętla przepisuje pixele o znalezionych kolorach
            //
            for (let sy = 0; sy < 8; sy++) {
              for (let sx = 0; sx < 8; sx++) {
                for (let k = 0; k < 3; k++) {
                  let pixelColor =
                    pixelsGLInverted.data[
                      (by + sy) * 320 * 4 + bx * 8 * 4 + sx * 4 + k
                    ];

                  if (!foundBackColor) {
                    if (pixelColor == newColor1 || pixelColor == newColor2) {
                      pixelsHires.data[
                        (by + sy) * 320 * 4 + bx * 8 * 4 + sx * 4 + k
                      ] = pixelColor;
                    }
                  } else {
                    if (pixelColor == newColor1) {
                      // console.log(
                      //   "nowy punkt [ " +
                      //     (by + sy) * 320 * 4 +
                      //     bx * 8 * 4 +
                      //     sx * 4 +
                      //     k +
                      //     " ] = " +
                      //     pixelColor
                      // );
                      pixelsHires.data[
                        (by + sy) * 320 * 4 + bx * 8 * 4 + sx * 4 + k
                      ] = pixelColor;
                    }
                  }
                }
              }
            }
          }
        }
        // console.log("pixelsHires:");
        // for (let i = 0; i < 320 * 4; i++) {
        //   console.log(pixelsHires[i]);
        // }
      }

      // Pokazanie róznic
      // ------------------------------------------------------------------------------------------
      for (let j = 0; j < 200; j++) {
        for (let i = 0; i < 320; i++) {
          for (let k = 0; k < 3; k++) {
            pixelsDiff.data[j * 320 * 4 + i * 4 + k] =
              pixelsGLInverted.data[j * 320 * 4 + i * 4 + k] -
              pixelsHires.data[j * 320 * 4 + i * 4 + k];
          }
        }
      }

      // test red dot
      // pixelsHires.data[320*4*100+160*4+0] = 0xff
      // pixelsHires.data[320*4*100+160*4+1] = 0
      // pixelsHires.data[320*4*100+160*4+2] = 0
      // pixelsHires.data[320*4*100+160*4+3] = 0xff

      // Zapis do output
      // ------------------------------------------------------------------------------------------
      contextHires.putImageData(pixelsHires, 0, 0);
      contextDiff.putImageData(pixelsDiff, 0, 0);

      // Zmiana background color na podstawie kliknięcia myszką
      // ------------------------------------------------------------------------------------------
      if (changeBackColor) {
        this.backgroundColorBack = this.colorString(
          backgroundColorRGBA[0],
          backgroundColorRGBA[1],
          backgroundColorRGBA[2],
          backgroundColorRGBA[3]
        );
        if (
          (backgroundColorRGBA[0] +
            backgroundColorRGBA[0] +
            backgroundColorRGBA[0]) /
            3 <
          0x80
        ) {
          this.backgroundColorText = this.colorString(0xff, 0xff, 0xff, 0xff);
        } else {
          this.backgroundColorText = this.colorString(0, 0, 0, 0xff);
        }
      }

      if (changeBackColor) {
        changeBackColor = false;
      }
    }
  },
  // ================================================================================================
  // Funkcje VUE
  // ================================================================================================
  mounted() {
    console.log("Mounted()");
    this.init();
    this.animate();
  }
};
</script>
