---
title:          Curse of the !DOCTYPE
date:           2016-02-01
author:         Sid Shukla 
authorLink:     https://github.com/thunderboltsid
description:    About the !DOCTYPE
---

DOCTYPE: The enigmatic syntactical declaration that’s forced down the throat of amateur web developers and designers without much understanding of what it implies. Sure, some person might’ve mentioned that it means document type, which is true, in a way, it does not quite reflect it’s true idea.

Presence or Absence of a DOCTYPE can lead to interesting and unpredictable (but not unforeseen) results. This started off at my previous internship. I was working on a visualization intensive web application that was supposed to help the company with better decision making. This was a typical django stack  with Angular.JS for the frontend, Postgres for database needs, and memcached for caching. One fine day, while refactoring the home page into a separate base layer  comprising of only the staticfiles (CSS + JS) which other pages can import, I broke the whole application.

Well, to say that would be a bit of an understatement, because random things stopped working that shouldn’t have had any connection with the changes or any reason to break in the first place. Hours (incl. Döner, Coffee, and Beer) were spent frustratingly trying to hunt down the origin of this bug . None of the resultant bugs made any sense with regards to the commit history. All of a sudden, things just stopped working and it was like a spell of voodoo. After frustrating hours, we just started thinking backwards. After a while, my colleague, who’s arguably way more experienced in the art of web dev than I am, started laughing uncontrollably as if he figured out everything and it turned out to be trivial. He asked me what was my “Brilliant idea of the day” and I remembered the DOCTYPE declaration that I added besides the refactoring. But, the thought was quickly brushed aside since I wondered it’s just a DOCTYPE. Something so basic shouldn’t cause a problem. I mean, wasn’t my browser adding it from me before, anyway? After I couldn’t find any other culprit, I decided to give it a shot and removed the DOCTYPE declaration.

Voila! Everything started working the way it was supposed to without any further modifications. This left me in a state of disbelief. How can something so benign, have such a significant impact on the state of an application. Now, if I were using some old school DOCTYPE, that might make a little bit of sense, but, no, I was using the regular, good old, HTML5 doctype. Whoops! What went wrong here?

To this day, it stays a mystery to me. Haunting me over those drunk escapades at times. Something that my puny human brain can’t possibly understand. We’re living in a world of over-smart browsers and perpetually stumped developers. One day, I’ll get to the bottom of this… Until then, Browsers – 1, Me – 0.

<iframe src="//giphy.com/embed/11kEuHSQAXXiGQ?html5=true" width="480" height="270" frameborder="0"></iframe>
