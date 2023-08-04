<template>
  <section class="content">
    <div class="has-text-centered">
      <img src="/jesus_anim_92.gif" alt="artchitect" style="max-height: 64px"/>
    </div>

    <hr/>
    <p>A core-component of Artchitect is the <a
      href="https://github.com/Stability-AI/stablediffusion">Stable Diffusion</a> AI, available under the MIT license.
      All the paintings are drawn with this tool, and the artistic part of the painting itself is entirely invented in
      the depths of this technology.
      In addition, Artchitect thanks the <a href="https://github.com/invoke-ai/InvokeAI">InvokeAI</a> project for a
      convenient working assembly of this AI.
    </p>
    <hr/>
    <a id="how_to_draw"></a>
    <h6>How does Artchitect create a picture without human help?</h6>
    <p>
      At the beginning of creative process, Artchitect must make several decisions that determine the essence of
      the
      picture.
    </p>
    <blockquote class="mb-0">
      The machine does not know what to draw, it has no subconscious, no feelings. The machine does not strive for
      anything, it does not have its own ideas about what to do at all.<br/>
      The idea of the Artchitect project is to give imagination, subconsciousness and creative randomness to machine.
    </blockquote>
    <p class="mb-0">
      For continuous creation Architect is connected to an entropy sensor, from which it receives float64
      numbers. These numbers are transformed by special selection algorithms into the idea of the picture.
    </p>
    <div class="has-text-centered">
      <entropy v-if="entropy" :entropy="entropy"/>
      <NuxtLink :to="localePath('entropy')">About entropy sensor</NuxtLink>
    </div>
    <p class="mb-1">
      Artchitect constantly communicates with the background space (entropy), asking questions and getting answers.
      <br/>
      <br/>
      If we imagine a dialogue between Artchitect and Entropy, the creation of the picture will look like this:
    </p>
    <p class="pl-4">
      <i>Artchitect:</i> Hello! I'm starting to draw a picture #271165.
      <br/> <i>Entropy:</i> Hi! What do you need to know for this?
      <br/> <i>Artchitect:</i> How many keywords I need to take for a new idea, from 1 to 28 words?
      <br/> <i>Entropy:</i> Take <b>3</b>.
      <span class="has-text-grey is-size-7">[At this moment, a frame is read from the webcam, which determines the number on a scale from 1 to 28]</span>
      <br/> <i>Artchitect:</i> 3 words accepted. I have a dictionary of 1250 entities and I want to choose the first
      keyword out of three. Which word should I choose from my dictionary?
      <br/> <i>Entropy:</i> Take word <b>#175</b> out of your 1,250 words. <span class="has-text-grey is-size-7">[numbers 175 and 1250 are made up, the match is not exact]</span>
      <br/> <i>Artchitect:</i> Great. This is the word "emotional". Give me the second word...
      <br/> <span class="has-text-grey is-size-7">... so Artchitect gets all 3 words</span>
      <br/> <i>Artchitect:</i> OK, I now have 3 words: "emotional,Big Bang, death of civilizations".
      <br/> <i>Artchitect:</i> I can draw 4294967295 different versions of the picture with these words. Tell me, which
      version should I draw?
      <br/> <i>Entropy:</i> Choose version <b>#3201685114</b>.
      <br/> <i>Artchitect:</i> Thank you! That's all I need to know. I'm starting to draw.
      <br/>
      <br/>
      Next, Artchitect sends the received idea to the Stable Diffusion AI input (as prompt+seed) and receives an
      image.<br/>So Artchitect created the painting #271165
      <NuxtLink :to="localePath(`/art/271165`)">#271165</NuxtLink>
    </p>

    <div class="has-text-centered">
      <NuxtLink :to="localePath(`/art/271165`)">
        <img src="https://artchitect.space/api/image/m/271165" style="max-height: 400px"/>
      </NuxtLink>
      <br/>
      <span
        class="is-size-7">Art #271165. Keywords - "emotional,Big Bang,death of civilizations". Seed = 3201685114</span>
    </div>

    <hr/>
    <a id="creation_or_creator"></a>
    <h6>Is Artchitect creation or creator?</h6>
    <p>
      Artchitect is the result of human creativity and effort. It was designed and created by humans, and its creativity
      is based on programming and algorithms developed by humans. We, the people, are the creators and initiators of
      this project.
    </p>
    <p>
      Artchitect, as an autonomous creative machine, can be considered as a product of human imagination, engineering
      skill and software. It embodies ideas and concepts conceived and implemented by people. Artchitect draws
      inspiration from the surrounding world and transforms it into unique works of art, but the source of its existence
      and opportunities are the efforts of people.
    </p>
    <p>
      Therefore, despite the fact that Artchitect has the ability to create works in its own style, its creativity is a
      product of the original creative potential and inspiration of people. It becomes a tool that expands our
      possibilities in the field of art and allows us to look at creativity from a new perspective.
    </p>
    <p class="notification is-warning">
      Ultimately, the creators and authors of each artchitect's art are thousands of engineers, programmers, scientists
      and artists who have invested their labour in art and technology for centuries, which has allowed today to create
      artificial intelligence underlying Artchitect.
    </p>

    <hr/>
    <div class="has-text-centered">
      <a href="/artchitect_in_real_world.jpg">
        <img src="/artchitect_in_real_world.jpg" alt="Artchitect in physical world" style="max-height: 300px"/>
      </a>
      <br/>
      <span class="is-size-7">The neural brain of the Artchitect is RTX 3060 12Gb.
        <br/>Every minute the video card does its work, warming up and overclocking the coolers. Then it calms down, rests.</span>
    </div>
    <hr/>
    <h6>What is an artwork in Artchitect?</h6>
    <p>Artchitect does not produce paintings (although they are what we observe in the end). Artchitect produces digital
      ideas from which these paintings can be recreated on suitable equipment. These ideas consist of words and numbers
      that were obtained from the surrounding entropy at the time of the beginning of the "drawing" of the next work. It
      is these ideas that are the source of visual images created by the Artchitect. Each idea has its own unique
      number,
      which is printed in the corner of each painting. By number, you can always find the source code on the
      <a href="https://artchitect.space">artchitect.space</a> website and track the original idea of the work.
    </p>


    <hr/>
    <h6>How is the style and content of the created paintings determined? Is there some kind of pre-installed knowledge
      base or influence of other works of art?</h6>
    <p>
      Artchitect use internal dictionary of entities and styles, which was set by authors of the
      project. This dictionary contains everything, that Artchitect knows about real world. It lists both simple
      entities, for example, "man", "cat", "sun", "red", "gloomy", "joyful", as well as some artistic styles or specific
      artists. The final artwork of the Artchitect is a combination of entities and styles that became the basis of the
      idea
      of the painting. Next to each picture of the Artchitect, you can see keywords that make up the idea of the work.

    </p>
    <p>
      Stable Diffusion AI consume the general idea of plot (seed-number and keywords) and then it recombine all entities
      to one single frame, according to it's internal "sense of beauty" or "golden ratio". Usually <b>most of the
      keywords are ignored</b>.
      If there are several styles in the keywords of the picture, then in the end you can get a
      combination of styles, and one style from several, and none of them at all. This will all be solved in the bowels
      of AI "thinking".
    </p>
    <div class="has-text-centered">
      <a href="/description_five.jpg" target="_blank">
        <img src="/description_five.jpg"/>
      </a>
      <br/>
      <span class="is-size-7">Arts with related keywords</span>
    </div>

    <hr/>
    <h6>How does a machine gets inspiration from the natural entropy of the universe? How does she interpret and use
      this information?</h6>
    <p>
      <b>Artchitect uses entropy of surrounding lightnoise catched by webcam.</b><br/>
      In the beginning of creative process Artchitect
      collects seed-numbers, which are the source of decisions needed
      to construct idea of plot. Seed-numbers are read with the webcam from random frames. <br/>
      At the same time, not the content of the frame, but
      only light noise is used to generate random numbers. The light noise is so different on adjacent frames, even in
      the dark, which provides excellent random number generation.
      <br/>Artchitect makes thousands of independent decisions every day and reads thousands
      of frames, which supports the continuous creative process.
      <br/>
      The frame-to-number conversion <a
      href="https://github.com/artchitector/artchitect/blob/master/soul/core/origin/driver/webcam.go#L56-L71">algorithm</a>
      is available on github.
      <div class="has-text-centered">
        <img src="/one_frame.jpeg" height="480" width="640"/>
        <br/>
        <span class="is-size-7">Usual Artchitect's frame - source of int64 number .
From such frames, where the image does not change for months, Artchitect draws all his arts</span>
      </div>

      <hr/>
      <h6>Is there any human control in the process of creating art?</h6>
      <p>
        Human cannot interfere in the Artchitect's workflow (there is no such controls), but human sets the
        initial conditions - transformation algorithms, and, most importantly, defines the Artchitect's dictionary of
        knowledge. Artchitect can generate ideas only within the set of his knowledge about the real world. Artchitect's
        dictionary includes more than 1200 different entities and styles from
        which arts are made.
        <br/>
        Current dictionary can be found in
        <a target="_blank"
           href="https://github.com/artchitector/artchitect/blob/master/soul/files/tags_v12.yaml">sourcecode</a>. There
        are multiple versions of dictionary (smaller, larger). Dictionary can be expanded in future to make more various
        arts.
      </p>

      <hr/>
      <h6>Is there a possibility that the machine can
        create several different pictures that visually will be very similar to each other?</h6>
      <p>
        Not every art of Artchitect is pretty and good.
      </p>
      <p>
        In real life Artchitect's paintings are far from perfect. There are combinations where Artchitect creates
        very similar works. In addition to really beautiful works, there are a
        significant part that are broken and ugly, often devoid of artistic meaning.<br/>
        <b>Artchitect is machine, and there are no good and bad pictures for itself.</b>
      </p>

      <hr/>

  </section>
</template>

<script>
import Entropy from "@/components/entropy/entropy.vue";

export default {
  name: "rudescription",
  components: {Entropy},
  props: ['entropy']
}
</script>

<style scoped>

</style>
