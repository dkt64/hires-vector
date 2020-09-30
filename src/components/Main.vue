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
    <v-btn @click="clickOnCalculate">Calculate...</v-btn>
    <v-row class="mt-4 ml-1 mb-4">Settings:</v-row>
    <v-chip
      class="mb-9"
      :text-color="backgroundColorText"
      :color="backgroundColorBack"
      >Background color (click on image)</v-chip
    >
    <!-- <v-slider
      style="width:500px"
      :thumb-size="24"
      thumb-label="always"
      v-model="max_sprites"
      step="1"
      min="1"
      max="8"
      label="Maximum sprites nr"
    >
      <template v-slot:append> -->
    <v-text-field
      v-model="rect_selection"
      class="mt-0 pt-0"
      type="number"
      style="width: 150px"
      label="Sprite rectangle selection"
    ></v-text-field>
    <v-text-field
      v-model="max_sprites"
      class="mt-0 pt-0"
      type="number"
      style="width: 150px"
      label="Max nr of sprites"
    ></v-text-field>
    <!-- </template>
    </v-slider> -->
    <v-switch class="mr-4" v-model="switch_controls" label="Animate"></v-switch>
    <v-switch class="mr-4" v-model="rysuj_siatke" label="Show grid"></v-switch>
    <v-text-field
      label="All colors"
      outlined
      :value="colors_all"
    ></v-text-field>
    <v-text-field
      label="Visible colors for hires bitmap"
      outlined
      :value="colors_bmp"
    ></v-text-field>
    <v-text-field
      label="Visible colors for sprites"
      outlined
      readonly
      :value="colors_spr"
    ></v-text-field>
    <v-text-field
      label="Sprite rectangles"
      outlined
      :value="rects_spr_txt"
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
    pixelsHires: null,
    pixelsSprites: null,
    colors_all: [],
    colors_bmp: [],
    colors_spr: [],
    sprites: [],
    rects_spr: [],
    rect_selection: -1,
    rects_spr_txt: "",
    slider_x: 0,
    slider_y: 0,
    slider_z: 0,
    max_sprites: 7,
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
    // drawGrid - rysowanie siatki
    // --------------------------------------------------------------------------------------------
    drawGrid: function(pixelsHires, pixelsSprites) {
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
    },
    fillImageData: function(image, r, g, b, a) {
      for (let j = 0; j < 200; j++) {
        for (let i = 0; i < 320; i++) {
          image.data[j * 320 * 4 + i * 4 + 0] = r;
          image.data[j * 320 * 4 + i * 4 + 1] = g;
          image.data[j * 320 * 4 + i * 4 + 2] = b;
          image.data[j * 320 * 4 + i * 4 + 3] = a;
        }
      }
    },
    // --------------------------------------------------------------------------------------------
    // countColors - liczenie kolorów
    // --------------------------------------------------------------------------------------------
    countColors: function(pixelsHires, pixelsSprites) {
      var found = false;

      var colors_spr = [];
      for (let j = 0; j < 200; j++) {
        for (let i = 0; i < 320; i++) {
          for (let k = 0; k < this.colors_all.length; k++) {
            if (this.colors_all[k] == pixelsSprites.data[j * 320 * 4 + i * 4]) {
              found = false;
              for (let l = 0; l < colors_spr.length; l++) {
                if (colors_spr[l] == this.colors_all[k]) {
                  found = true;
                }
              }
              if (!found) {
                colors_spr.push(this.colors_all[k]);
              }
            }
          }
        }
      }
      this.colors_spr = colors_spr;

      var colors_bmp = [];
      for (let j = 0; j < 200; j++) {
        for (let i = 0; i < 320; i++) {
          for (let k = 0; k < this.colors_all.length; k++) {
            if (this.colors_all[k] == pixelsHires.data[j * 320 * 4 + i * 4]) {
              found = false;
              for (let l = 0; l < colors_bmp.length; l++) {
                if (colors_bmp[l] == this.colors_all[k]) {
                  found = true;
                }
              }
              if (!found) {
                colors_bmp.push(this.colors_all[k]);
              }
            }
          }
        }
      }
      this.colors_bmp = colors_bmp;
    },
    // --------------------------------------------------------------------------------------------
    // drawRects - rysowanie prostokątów
    // --------------------------------------------------------------------------------------------
    drawRects: function(image) {
      if (this.rects_spr != null) {
        var k1 = this.rect_selection;
        var k2;

        if (k1 >= 0 && k1 < this.rects_spr.length) {
          k2 = k1;
        } else {
          k1 = 0;
          k2 = this.rects_spr.length - 1;
        }

        if (this.rects_spr[k1] != null && this.rects_spr[k2] != null) {
          for (let k = k1; k <= k2; k++) {
            var x0 = this.rects_spr[k].x0;
            var y0 = this.rects_spr[k].y0;
            var x1 = this.rects_spr[k].x1;
            var y1 = this.rects_spr[k].y1;

            // Pozioma górna
            var dy = y0;
            for (let dx = x0; dx <= x1; dx += 2) {
              // image.data[dy * 320 * 4 + dx * 4 + 0] = k * 0x10;
              // image.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
              image.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
              // image.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
            }

            // Pozioma dolna
            dy = y1;
            for (let dx = x0; dx <= x1; dx += 2) {
              // image.data[dy * 320 * 4 + dx * 4 + 0] = 0x80;
              // image.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
              image.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
              // image.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
            }

            // Pionowa lewa
            var dx = x0;
            for (let dy = y0; dy <= y1; dy += 2) {
              // image.data[dy * 320 * 4 + dx * 4 + 0] = 0x80;
              // image.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
              image.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
              // image.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
            }

            // Pionowa prawa
            dx = x1;
            for (let dy = y0; dy <= y1; dy += 2) {
              // image.data[dy * 320 * 4 + dx * 4 + 0] = 0x80;
              // image.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
              image.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
              // image.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
            }
          }
        }
      }
    },
    // --------------------------------------------------------------------------------------------
    // findRects - znaleienie prostokątów
    // --------------------------------------------------------------------------------------------
    findRects: function(pixelsSprites) {
      var rects_spr = [];

      // Pętla z kolorami
      for (let k = 0; k < this.colors_spr.length; k++) {
        // Selekcja plane'a
        var plane = new ImageData(320, 200);
        this.fillImageData(plane, 0, 0, 0, 0xff);
        for (let j = 0; j < 200; j++) {
          for (let i = 0; i < 320; i++) {
            if (this.colors_spr[k] == pixelsSprites.data[j * 320 * 4 + i * 4]) {
              plane.data[j * 320 * 4 + i * 4] =
                pixelsSprites.data[j * 320 * 4 + i * 4];
              // console.log("found " + k);
            }
          }
        }
        // Szukamy x0,y0,x1,y1 (x_min, x_max, y_min, y_max)
        var x0 = 0;
        var y0 = 0;
        var x1 = 319;
        var y1 = 199;

        // y0
        for (let j = 0; j < 200; j++) {
          for (let i = 0; i < 320; i++) {
            if (plane.data[j * 320 * 4 + i * 4] != 0) {
              y0 = j;
              j = 200;
              break;
            }
          }
        }
        // y1
        for (let j = 199; j >= 0; j--) {
          for (let i = 0; i < 320; i++) {
            if (plane.data[j * 320 * 4 + i * 4] != 0) {
              y1 = j;
              j = -1;
              break;
            }
          }
        }
        // x0
        for (let i = 0; i < 320; i++) {
          for (let j = 0; j < 200; j++) {
            if (plane.data[j * 320 * 4 + i * 4] != 0) {
              x0 = i;
              i = 320;
              break;
            }
          }
        }
        // x1
        for (let i = 319; i >= 0; i--) {
          for (let j = 0; j < 200; j++) {
            if (plane.data[j * 320 * 4 + i * 4] != 0) {
              x1 = i;
              i = -1;
              break;
            }
          }
        }

        // Dodajemy nowy rect
        var newRect = {
          i: k,
          x0: x0,
          x1: x1,
          y0: y0,
          y1: y1,
        };

        rects_spr.push(newRect);
        // this.planes.push(plane)
      }

      // Kopiowanie do zm. globalnej
      this.rects_spr = rects_spr;
      console.log(this.rects_spr);

      // Wygenerowanie opisu
      var rects_spr_txt = "";
      for (let i = 0; i < rects_spr.length; i++) {
        rects_spr_txt += rects_spr[i].i + ": ";
        rects_spr_txt += "x0=" + rects_spr[i].x0;
        rects_spr_txt += " x1=" + rects_spr[i].x1;
        rects_spr_txt += " y0=" + rects_spr[i].y0;
        rects_spr_txt += " y1=" + rects_spr[i].y1;
        rects_spr_txt += "    ";
      }
      this.rects_spr_txt = rects_spr_txt;
    },
    // --------------------------------------------------------------------------------------------
    // findSprites - wygenerowanie sprajtów z prostokątów
    // --------------------------------------------------------------------------------------------
    // findSprites: function(pixelsSprites) {
      // if (this.rects_spr != null) {
      //   for (let i = 0; i < this.rects_spr.length; i++) {
      //     var x0 = this.rects_spr[i].x0;
      //     var y0 = this.rects_spr[i].y0;
      //     var x1 = this.rects_spr[i].x1;
      //     var y1 = this.rects_spr[i].y1;

      //     var w_prawo
      //     // Sprawdzamy kierunek
      //     //
      //     if (image.data[y0 * 320 * 4 + x0 * 4 + 0] != 0) {
      //       w_prawo = true
      //     } 
      //     if (image.data[y0 * 320 * 4 + x1 * 4 + 0] != 0) {
      //       w_prawo = false
      //     } 

      //     if (w_prawo) {

      //     }

      //     for (let dx = x0; dx <= x1; dx++) {
      //       // image.data[dy * 320 * 4 + dx * 4 + 0] = k * 0x10;
      //       // image.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
      //       image.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
      //       // image.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
      //     }


      //     var spr_x = 0;
      //     var spr_y = 0;

      //     // Dodajemy nowy rect
      //     var newSprite = {
      //       color_nr: i,
      //       poz_x: spr_x,
      //       poz_y: spr_y,
      //       data: [24 * 21],
      //     };


      //     // Pozioma dolna
      //     dy = y1;
      //     for (let dx = x0; dx <= x1; dx += 2) {
      //       // image.data[dy * 320 * 4 + dx * 4 + 0] = 0x80;
      //       // image.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
      //       image.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
      //       // image.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
      //     }

      //     // Pionowa lewa
      //     var dx = x0;
      //     for (let dy = y0; dy <= y1; dy += 2) {
      //       // image.data[dy * 320 * 4 + dx * 4 + 0] = 0x80;
      //       // image.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
      //       image.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
      //       // image.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
      //     }

      //     // Pionowa prawa
      //     dx = x1;
      //     for (let dy = y0; dy <= y1; dy += 2) {
      //       // image.data[dy * 320 * 4 + dx * 4 + 0] = 0x80;
      //       // image.data[dy * 320 * 4 + dx * 4 + 1] = 0x80;
      //       image.data[dy * 320 * 4 + dx * 4 + 2] = 0x80;
      //       // image.data[dy * 320 * 4 + dx * 4 + 3] = 0xff;
      //     }
      //   }
      // }
    // },
    // --------------------------------------------------------------------------------------------
    // Sprites - wygenerowanie sprajtów
    // --------------------------------------------------------------------------------------------
    Sprites: function(pixelsSprites) {
      this.findRects(pixelsSprites);
      // this.findSprites(pixelsSprites);
    },
    // --------------------------------------------------------------------------------------------
    // clickOnCalculate - kliknięcie na context GL
    // --------------------------------------------------------------------------------------------
    clickOnCalculate: function() {
      this.countColors(this.pixelsHires, this.pixelsSprites);
      this.Sprites(this.pixelsSprites);
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

      var colorStep = 200 / this.izoGeo.faces.length;
      for (let i = 0; i < this.izoGeo.faces.length; i++) {
        //izoGeo.faces[i].color.setHex((c64colors[i & 0x0f] * (i + 1)) & 0xffff00 + i);
        this.izoGeo.faces[i].color.setHex(
          0x202020 + i * (colorStep + colorStep * 256 + colorStep * 256 * 256)
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
      this.fillImageData(pixelsHires, 0, 0, 0, 0xff);
      this.fillImageData(pixelsSprites, 0, 0, 0, 0xff);

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

      // Malowanie siatki na obrazie
      // ------------------------------------------------------------------------------------------
      this.drawGrid(pixelsHires, pixelsSprites);

      // Malowanie siatki na obrazie
      // ------------------------------------------------------------------------------------------
      this.drawRects(pixelsSprites);

      // Liczenie kolorów
      // ------------------------------------------------------------------------------------------
      // this.countColors(pixelsHires, pixelsSprites);

      // Analiza obrazu - Sprites
      // ------------------------------------------------------------------------------------------
      // this.Sprites(pixelsSprites);

      // Zapis do output
      // ------------------------------------------------------------------------------------------
      contextHires.putImageData(pixelsHires, 0, 0);
      contextSprites.putImageData(pixelsSprites, 0, 0);
      this.pixelsHires = pixelsHires;
      this.pixelsSprites = pixelsSprites;

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
