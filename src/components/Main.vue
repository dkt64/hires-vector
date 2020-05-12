<template>
  <v-container>
    <v-row class="d-flex justify-space-around">
      <v-card tile outlined class="mx-auto pa-2">
        <v-toolbar color="indigo" class="headline" dark>Three.js</v-toolbar>
        <div
          v-on:click="clickOnGL"
          id="canvasGL"
          style="width: 320px; height: 200px"
        ></div>
      </v-card>
      <v-card tile outlined class="mx-auto pa-2">
        <v-toolbar color="indigo" class="headline" dark>Hires</v-toolbar>
        <canvas id="canvasHires" width="320" height="200"></canvas>
      </v-card>
      <v-card tile outlined class="mx-auto pa-2">
        <v-toolbar color="indigo" class="headline" dark>Sprites</v-toolbar>
        <canvas id="canvasSprites" width="320" height="200"></canvas>
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
    <v-divider class="mt-2 mb-4"></v-divider>
    <v-row class="ml-1 mb-4">Settings:</v-row>
    <v-chip
      class="mb-9"
      :text-color="backgroundColorText"
      :color="backgroundColorBack"
      >Background color (click on image)</v-chip
    >
    <v-slider
      style="width:500px"
      :thumb-size="24"
      thumb-label="always"
      v-model="slider_sprites_nr"
      step="1"
      min="1"
      max="8"
      label="Maximum sprites nr"
    >
      <template v-slot:append>
        <v-text-field
          v-model="slider_sprites_nr"
          class="mt-0 pt-0"
          hide-details
          single-line
          type="number"
          style="width: 60px"
        ></v-text-field>
      </template>
    </v-slider>
    <v-switch class="mr-4" v-model="switch_controls" label="Animate"></v-switch>
    <v-switch class="mr-4" v-model="rysuj_siatke" label="Show grid"></v-switch>
    <v-text-field
      label="All colors"
      outlined
      readonly
      :value="colors_all"
    ></v-text-field>
    <v-text-field
      label="Visible colors for sprites"
      outlined
      readonly
      :value="colors_sprites"
    ></v-text-field>
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

export default {
  name: "Main",

  // ================================================================================================
  // Zmienne VUE
  // ================================================================================================
  data: () => ({
    renderer: null,
    scene: null,
    camera: null,
    izoGeo: null,
    izoMat: null,
    izo: null,
    colors_all: [],
    colors_sprites: [],
    slider_x: 0,
    slider_y: 0,
    slider_z: 0,
    slider_sprites_nr: 0,
    switch_controls: false,
    backgroundColorBack: "#000000ff",
    backgroundColorText: "#ffffffff",
    rysuj_siatke: false,
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
        45,
        canvasGL.clientWidth / canvasGL.clientHeight,
        0.1,
        1000
      );
      this.camera.position.z = 200;

      this.scene = new THREE.Scene();

      const light = new THREE.AmbientLight(0xffffff);
      // light.position.set(0, 0, 1000);
      this.scene.add(light);

      this.izoGeo = new THREE.IcosahedronGeometry(70);
      this.izoMat = new THREE.MeshPhongMaterial({
        color: 0xffffff,
        vertexColors: THREE.FaceColors,
      });

      for (let i = 0; i < this.izoGeo.faces.length; i++) {
        //izoGeo.faces[i].color.setHex((c64colors[i & 0x0f] * (i + 1)) & 0xffff00 + i);
        this.izoGeo.faces[i].color.setHex(
          0x202020 + (i * 8 + 1) + (i * 8 + 1) * 256 + (i * 8 + 1) * 256 * 256
        );
      }

      var colors_all = new Uint8Array(this.izoGeo.faces.length);
      for (let i = 0; i < this.izoGeo.faces.length; i++) {
        colors_all[i] = this.izoGeo.faces[i].color.getHex();
      }
      this.colors_all = colors_all;
      console.log("Colors:");
      console.log(this.colors);

      this.izo = new THREE.Mesh(this.izoGeo, this.izoMat);
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

      var canvasSprites = document.getElementById("canvasSprites");
      var contextSprites = canvasSprites.getContext("2d");

      // Przygotowanie buforów na pixele
      // ------------------------------------------------------------------------------------------
      // var pixelsHires = contextHires.getImageData(0, 0, 320, 200);
      var pixelsGLInverted = new ImageData(320, 200);
      var pixelsHires = new ImageData(320, 200);
      var pixelsSprites = new ImageData(320, 200);

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
      // if (changeBackColor) {
      // console.log("pixelsGLInverted:");
      // console.log(pixelsGLInverted);

      // Czyścimy bufory
      //
      // ------------------------------------------------------------------------------------------
      for (let j = 0; j < 200; j++) {
        for (let i = 0; i < 320; i++) {
          pixelsHires.data[j * 320 * 4 + i * 4 + 0] = 0;
          pixelsHires.data[j * 320 * 4 + i * 4 + 1] = 0;
          pixelsHires.data[j * 320 * 4 + i * 4 + 2] = 0;
          pixelsHires.data[j * 320 * 4 + i * 4 + 3] = 0xff;
          pixelsSprites.data[j * 320 * 4 + i * 4 + 0] = 0;
          pixelsSprites.data[j * 320 * 4 + i * 4 + 1] = 0;
          pixelsSprites.data[j * 320 * 4 + i * 4 + 2] = 0;
          pixelsSprites.data[j * 320 * 4 + i * 4 + 3] = 0xff;
        }
      }

      // Szukamy obszarów 8x8
      //
      for (let by = 0; by < 200; by += 8) {
        for (let bx = 0; bx < 320; bx += 8) {
          var foundBackColor = false;
          var newColor1 = backgroundColorRGBA[0];
          var newColor2 = backgroundColorRGBA[0];

          // Pierwsza pętla szuka background kolor
          // Porównujemy tylko pierwszy bajt z RGBA
          //
          for (let sy = 0; sy < 8; sy++) {
            for (let sx = 0; sx < 8; sx++) {
              // Odczyt pixela
              //
              let pixelColor =
                pixelsGLInverted.data[(by + sy) * 320 * 4 + (bx + sx) * 4];

              // Sprawdzamy czy background
              // Jeżeli tak to zaznaczamy że znaleźliśmy
              //
              if (pixelColor == backgroundColorRGBA[0]) {
                foundBackColor = true;
              }
            }
          }

          // Druga pętla szuka kolorów dla Hires
          // Porównujemy tylko pierwszy bajt z RGBA
          //
          for (let sy = 0; sy < 8; sy++) {
            for (let sx = 0; sx < 8; sx++) {
              // Odczyt pixela
              //
              let pixelColor =
                pixelsGLInverted.data[(by + sy) * 320 * 4 + (bx + sx) * 4];

              // Kolor 1
              //
              if (
                newColor1 == backgroundColorRGBA[0] &&
                pixelColor != newColor1 &&
                pixelColor != backgroundColorRGBA[0]
              ) {
                newColor1 = pixelColor;
              }

              // Kolor 2
              //
              if (!foundBackColor) {
                if (
                  newColor2 == backgroundColorRGBA[0] &&
                  pixelColor != newColor2 &&
                  pixelColor != newColor1 &&
                  pixelColor != backgroundColorRGBA[0]
                ) {
                  newColor2 = pixelColor;
                }
              }
            }
          }

          // Trzecia pętla przepisuje pixele o znalezionych kolorach
          //
          for (let sy = 0; sy < 8; sy++) {
            for (let sx = 0; sx < 8; sx++) {
              let pixelColor =
                pixelsGLInverted.data[(by + sy) * 320 * 4 + (bx + sx) * 4];

              // Gdy w kwadracie 8x8 znajduje się kolor tła to możemy uzyć tylko jednego koloru
              //
              if (!foundBackColor) {
                if (pixelColor == newColor1 || pixelColor == newColor2) {
                  pixelsHires.data[
                    (by + sy) * 320 * 4 + (bx + sx) * 4
                  ] = pixelColor;
                }
              } else {
                if (pixelColor == newColor1) {
                  pixelsHires.data[
                    (by + sy) * 320 * 4 + (bx + sx) * 4
                  ] = pixelColor;
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
      // }

      // Wypełniamy kanał ALFA
      // przepisujemy bajt 0 na 1+2
      // ------------------------------------------------------------------------------------------
      for (let j = 0; j < 200; j++) {
        for (let i = 0; i < 320; i++) {
          pixelsHires.data[j * 320 * 4 + i * 4 + 1] =
            pixelsHires.data[j * 320 * 4 + i * 4 + 0];
          pixelsHires.data[j * 320 * 4 + i * 4 + 2] =
            pixelsHires.data[j * 320 * 4 + i * 4 + 0];
          pixelsHires.data[j * 320 * 4 + i * 4 + 3] = 0xff;
          pixelsSprites.data[j * 320 * 4 + i * 4 + 3] = 0xff;
        }
      }

      // Pokazanie róznic
      // ------------------------------------------------------------------------------------------
      for (let j = 0; j < 200; j++) {
        for (let i = 0; i < 320; i++) {
          for (let k = 0; k < 3; k++) {
            pixelsSprites.data[j * 320 * 4 + i * 4 + k] =
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

      // Rysujemy linie kropkowane
      // pionowe
      // ------------------------------------------------------------------------------------------
      if (this.rysuj_siatke) {
        for (let dy = 0; dy < 200; dy += 2) {
          for (let dx = 0; dx < 320; dx += 8) {
            pixelsHires.data[dy * 320 * 4 + dx * 4 + 0] = 0x80;
            pixelsHires.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
            pixelsHires.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
            pixelsHires.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
            pixelsSprites.data[dy * 320 * 4 + dx * 4 + 0] = 0x80;
            pixelsSprites.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
            pixelsSprites.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
            pixelsSprites.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
          }
        }

        // Rysujemy linie kropkowane
        // poziome
        // ------------------------------------------------------------------------------------------
        for (let dy = 0; dy < 200; dy += 8) {
          for (let dx = 0; dx < 320; dx += 2) {
            pixelsHires.data[dy * 320 * 4 + dx * 4 + 0] = 0x80;
            pixelsHires.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
            pixelsHires.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
            pixelsHires.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
            pixelsSprites.data[dy * 320 * 4 + dx * 4 + 0] = 0x80;
            pixelsSprites.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
            pixelsSprites.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
            pixelsSprites.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
          }
        }
      }

      // Kolory
      // ------------------------------------------------------------------------------------------
      // var colors_sprites = new Uint8Array();
      var colors_sprites = []
      for (let j = 0; j < 200; j++) {
        for (let i = 0; i < 320; i++) {
          for (let k = 0; k < this.colors_all.length; k++) {
            if (this.colors_all[k] == pixelsSprites.data[j * 320 * 4 + i * 4]) {
              var found = false;
              for (let l = 0; l < colors_sprites.length; l++) {
                if (colors_sprites[l] == this.colors_all[k]) {
                  found = true;
                }
              }
              if (!found) {
                colors_sprites.push(this.colors_all[k]);
              }
            }
          }
        }
      }
      this.colors_sprites = colors_sprites;

      // Zapis do output
      // ------------------------------------------------------------------------------------------
      contextHires.putImageData(pixelsHires, 0, 0);
      contextSprites.putImageData(pixelsSprites, 0, 0);

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
    },
  },
  // ================================================================================================
  // Funkcje VUE
  // ================================================================================================
  mounted() {
    console.log("Mounted()");
    this.init();
    this.animate();
  },
};
</script>
