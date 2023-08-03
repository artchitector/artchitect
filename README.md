# Artchitect version2 (re-engineering of v1)

![artchitect_logo](https://github.com/artchitector/artchitect2/blob/main/files/images/logo_anim_92.gif)

#### https://artchitect.space

> Artchitect - it is an amazing autonomous creative machine capable of creating magnificent artworks inspired by the
> universe around us. In its continuous creativity , the machine receives inspiration from the natural entropy of the
> universe, represented as background light, and creates unique artworks without human participation.

Techically artchitect-project is the control-system wrapped around an art-system - Stable Diffusion AI v1.5. Stable
Diffusion AI is the ability to draw pictures for Artchitect.
The data to run Architect is light background/noise. The light background is read using a webcam, the frame is converted
to an int64 number. **Int64** is the source of all solutions (randomly select something from the list or create a unique
initial value).

Arts examples:
![examples](https://github.com/artchitector/artchitect2/blob/main/files/images/description_five.jpg)

### architecture of artchitect

Two parts:

- soul - home computer with strong GPU to run Stable
  Diffusion (artchitect using RTX 3060 12Gb and Stable Diffusion v1.5, Invoke.AI "fork").
- gate - multiple dedicated VDSs (frontend+backend server, file storage, database)

Home computer need access to VDS, but VDS doesn't need access to home computed.

More GPU-RAM = larger resolution = more quality. RTX3060-12Gb gives resolution 2560x3840 (printable on 40x60 canvas).

golang backend services + python backend services, splitted between home computer and remote VDS (visible from
Internet).

### manual

🤝 there are no instructions for Artchitect, since no one needed it before.
If you need more information or instructions to install your copy of Artchitect - please ask me questions in issues. I
will help.

### engineering style

- no restrictions. artchitect is available for everyone. If you want your running-copy of
  Artchitect, you can do it. (
  But you need devops skills to
  setup and run your own servers, databases, to understand code)

```
/*
* — * «THE BEER-WARE LICENSE» (Revision 42):
* <i@nkuhta.ru> wrote this file. As long as you retain this notice you
* can do whatever you want with this stuff. If we meet some day, and you think
* this stuff is worth it, you can buy me a beer in return. Poul-Henning Kamp
* — */
```

- no reliability. SLA is not important. If the Artchitect breaks down and turns off, it can stand for
  several days
  before repair. VDS servers that provide viewing of paintings work reliably, but the "soul" - main home computer can
  be turned off for downtime from days to weeks if necessary.

- less seriously, more jokes. many things in source-code are called mystical and not obvious: main
  service called "soul", and it consists of "gifter", "merciful", "speller"... (sometimes the author does not remember
  which of them does what). `Don't be boring while code review.`

- no docs. if you need help, make an issue. author will personally explain what to do or will make docs

- no braches, everything is in main-branch in singlerepo. Almost every commit is delivered to production.

- no testing, neither manual nor automated

- no deployment automation, manual ssh-commands and crontasks

### How Artchitect looks like:

![artchitect_installation](https://github.com/artchitector/artchitect2/blob/main/files/images/artchitect_hardware.jpg)