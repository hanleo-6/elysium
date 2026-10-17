  ## Elysium

  > A terminal-rendered leaf sea slug ecosystem written in Go.



https://github.com/user-attachments/assets/46d9f516-5c88-4ee1-bd14-f2093f164092




  ### Concept

  One of my favourite seaslugs is *Elysia chlorotica*.
  
  **Elysium** is a simulation centred on their biology.
  
  With some artistic licence, this is the story of the simulation:

  As juveniles, the seaslugs greedily seek, attack* and nom algae. You might
  know that algae contain chloroplasts. Much like the first multi-cellular
  organism, these seaslugs decide to not digest these, instead they store them
  within their digestive system.

  After feeding on enough algae, the seaslugs undergo **metamorphosis**, gradually
  becoming green in the process.
  They can now harness the power of photosynthesis through the chloroplasts
  they have incorporated into their own cells!
  At this point, their parapodia unfurl, giving them their distinctive and
  beautiful leaf shape.

  Since their hustle and bustle of their young days is now over, its time for
  the adult seaslugs to reep the rewards of their hard work, as many of us
  dream of.
  They do this by **loafing**.

  In the process of loafing, they will seek out the underwater shafts of
  light, then proceed to loaf in them.

  When they've enjoyed their time enough, they will call down a **wormhole**
  to go through and be reborn as a juvenile, repeating their life cycle
  until you, as 'the person upstairs' presses Ctrl-C in the terminal.

  

  *seaslugs puncture the 'algal cell wall with its radula'. how violent!

---

  ### Features

  - **Procedural animation::** Seaslugs are modelled with a chain of point
  masses with distance constraints between them. The head seek-steers
  (Reynolds algorithm) towards nearest algae. Sprites are rendered
  dynamically using the point mass locations, resulting in a look inspired by
  'Rainworld'.
  - **Unicode Braille dot canvas:** Rendered in the terminal using 2×4 Braille
  cell dot matrices with true 24-bit ANSI color and cell background shading. Note: a limitation of using Braille dots is that
  each set of dots in a cell must be the same colour.
  - **Aquatic Hydrodynamic Fluid Dynamic Provider Abstract Physics Factory:** Algae flow with a vector field
  driven by curl noise, and are pushed aside by each seaslug's wake (like a boat wake).
  - **Underwater lighting:** A dynamic underwater light cone sways mimicking sunlight refraction. The maths is
  3 sin waves of different phases with 'Ornstein–Uhlenbeck' centre correction to pull it back,

---

  ### Running
  #### Requirements
  - Go 1.26+ (or compatible modern Go toolchain)
  - A terminal emulator supporting **24-bit TrueColor** (I currently use
  Alacritty, but Kitty, iTerm2, WezTerm, etc. should work too)
  and UTF-8 Braille characters.

  #### Run
  ```sh
  go run ./cmd/elysium
  ```

  Quit: `q` or `Ctrl+C`
  
---

  ### Architecture

  ```
  elysium/
  ├── cmd/
  │   └── elysium/main.go          # bubbletea runtime, simulation loop
  └── internal/
      ├── physics/                 # physics primitives and point mass structs
      ├── sim/                     # ecosystem logic: slugs, algae, curl flow & wake grid, saltmarsh lighting
      └── render/                  # rendering logic for each component of the scene
  ```

---

  ### Dependencies

  - **Terminal TUI** | [`charm.land/bubbletea/v2`](https://github.com/charmbracelet/bubbletea)
  - **Color Spaces** | [`github.com/lucasb-eyer/go-colorful`](https://github.com/lucasb-eyer/go-colorful)

 ### Inspirations

  - **Rainworld** | [`rainworldgame.com/`](https://rainworldgame.com/)
  - **cbonsai** | [`github.com/jakobrees/cbonsai`](https://github.com/jakobrees/cbonsai)
  - **asciiquarium** | [`github.com/antoineFrau/asciiquarium/tree/main`](https://github.com/antoineFrau/asciiquarium/tree/main)
  - **drawille** | [`github.com/asciimoo/drawille/tree/master`](https://github.com/asciimoo/drawille/tree/master)

---

 ### Current next steps!

  - ~~write list of next steps~~
  - rewrite lighting to tint foreground colours
  - improve parapodia rendering and animation
  - write full lifecycle logic for seaslugs 
  - implement wormholes
  - allow loafing
