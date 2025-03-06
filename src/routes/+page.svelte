<!-- Hero -->
<script>
  import "../app.css";
  import "../mybookmark.css";
  import { ModeWatcher, setMode } from "mode-watcher";
  import { onMount } from "svelte";
  import { Button } from "$lib/components/ui/button";
  import { Plus, ChevronDown, Github } from "lucide-svelte";
  import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "$lib/components/ui/card";

  onMount(() => {});

  setMode("dark");

  const features = [
    {
      title: "Improved Search",
      description:
        "Search not just url but content of the page. Get the best results effortlessly. No need to add labels or tags.",
      icon: "💾",
      image: "landing/sch.png"
    },
    {
      title: "Import Github Repos",
      description:
        "Import your starred repositories from Github with a single click. Always find the repo you starred long time ago",
      icon: "🔄",
      image: "landing/gh.png"
    },
    {
      title: "Import Bookmarks",
      description: "Upload bookmarks from your browsers and other bookmarking tools in a few clicks.",
      icon: "🔍",
      image: "landing/browser.png"
    },
    {
      title: "Schedules",
      description: "Add schedules to fetch the content from github trending or product hunt leaderboard",
      icon: "👥",
      image: "landing/schedules.png"
    }
  ];

  const pricingPlans = [
    {
      name: "Free",
      price: "$0",
      description: "Perfect for casual users",
      features: ["Up to 100 bookmarks", "Basic organization", "Standard support"],
      cta: "Start Free"
    },
    {
      name: "Pro",
      price: "$9.99/mo",
      description: "For power users and professionals",
      features: ["Unlimited bookmarks", "Advanced AI features", "Priority support", "Team sharing"],
      cta: "Get Pro",
      popular: true
    },
    {
      name: "Enterprise",
      price: "Custom",
      description: "For organizations with advanced needs",
      features: ["Everything in Pro", "Custom integrations", "Dedicated support", "Admin controls"],
      cta: "Contact Sales"
    }
  ];

  const testimonials = [
    {
      quote: "Smaraka transformed how I organize my research materials. I can't imagine working without it now!",
      author: "Sarah J.",
      role: "Content Researcher"
    },
    {
      quote: "The AI-powered categorization saves me hours of manual organization every week.",
      author: "Michael T.",
      role: "Data Scientist"
    },
    {
      quote: "Our team's productivity increased by 30% after adopting Smaraka for our shared resources.",
      author: "Elena R.",
      role: "Project Manager"
    }
  ];

  const faqs = [
    {
      question: "How does Smaraka differ from traditional bookmarking?",
      answer:
        "Smaraka uses AI to automatically categorize, tag, and make your bookmarks searchable in ways traditional tools can't. It's designed for knowledge workers who need smart organization."
    },
    {
      question: "Can I import my existing bookmarks?",
      answer:
        "Absolutely! Smaraka allows easy importing from all major browsers and bookmark managers with just a few clicks."
    },
    {
      question: "Is my data secure?",
      answer:
        "We use industry-leading encryption and security practices. Your data is stored securely and never shared with third parties."
    },
    {
      question: "How does the team sharing feature work?",
      answer:
        "Pro and Enterprise users can create shared collections that team members can access, contribute to, and collaborate on in real-time."
    }
  ];
</script>

<!-- Navbar -->
<header class="fixed top-0 z-50 w-full bg-background/80 backdrop-blur-md">
  <div class="container flex h-16 items-center justify-between">
    <div class="flex items-center gap-2">
      <img src="smaraka.png" alt="Smaraka Logo" class="h-8 w-auto" />
      <span class="text-2xl font-bold">Smaraka</span>
    </div>

    <div class="flex items-center gap-4">
      <Button variant="outline">
        <Github class="mr-2 h-4 w-4" />
        Fork me on Github
      </Button>
    </div>
  </div>
</header>

<!-- Hero Section -->
<section class="pb-12 pt-32">
  <div class="container flex flex-col items-center text-center">
    <div class="max-w-3xl space-y-4">
      <h1 class="text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl">Remember everything that matters</h1>
      <p class="mx-auto max-w-prose text-xl text-muted-foreground">
        Smaraka is an open source intelligent bookmarking platform that helps you organize and find digital content
        effortlessly.
      </p>
    </div>

    <div class="mt-10 flex flex-wrap justify-center gap-4">
      <Button size="lg">Documentation</Button>
    </div>

    <div
      class="elative mx-auto mt-16 w-full max-w-5xl overflow-hidden rounded-lg border border-foreground shadow-2xl ring ring-gray-200 ring-opacity-20 ring-offset-2 ring-offset-gray-400"
    >
      <img src="hero_1.png" alt="Smaraka dashboard" class="w-full" />
    </div>
  </div>
</section>

<!-- Features Section -->
<section id="features" class="bg-background py-24">
  <div class="container">
    <div class="mx-auto mb-16 max-w-2xl text-center">
      <h2 class="mb-4 text-3xl font-bold">Features that simplify your digital life</h2>
      <p class="text-muted-foreground">
        Designed to help you organize and access your digital content with powerful tools that adapt to your needs.
      </p>
    </div>

    <div class="mx-auto grid gap-16 gap-y-16 md:max-w-[875px] md:grid-cols-2">
      {#each features as feature, i}
        <div class="col-span-1 flex flex-col gap-4">
          <div
            class="flex-1 rounded-lg bg-background p-4 ring ring-gray-200 ring-opacity-20 ring-offset-2 ring-offset-gray-400"
          >
            <img src={feature.image} alt={feature.title} class=" w-full" />
          </div>
          <div class="flex flex-1 flex-col justify-center">
            <h3 class="mb-2 text-xl font-semibold">{feature.title}</h3>
            <p class="text-muted-foreground">{feature.description}</p>
          </div>
        </div>
      {/each}
    </div>
  </div>
</section>

<!-- Pricing Section -->

<!-- Testimonials Section -->
<section id="testimonials" class="bg-muted/50 py-24">
  <div class="container">
    <div class="mx-auto mb-16 max-w-2xl text-center">
      <h2 class="mb-4 text-3xl font-bold">What our users say</h2>
      <p class="text-muted-foreground">
        Thousands of professionals rely on Smaraka every day to organize their digital world.
      </p>
    </div>

    <div class="mx-auto grid max-w-5xl gap-8 md:grid-cols-3">
      {#each testimonials as testimonial}
        <Card>
          <CardContent class="pt-6">
            <div class="mb-4 text-2xl">❝</div>
            <p class="mb-6 italic">{testimonial.quote}</p>
            <div>
              <p class="font-semibold">{testimonial.author}</p>
              <p class="text-sm text-muted-foreground">{testimonial.role}</p>
            </div>
          </CardContent>
        </Card>
      {/each}
    </div>
  </div>
</section>

<!-- FAQ Section -->
<section id="faq" class="py-24">
  <div class="container">
    <div class="mx-auto mb-16 max-w-2xl text-center">
      <h2 class="mb-4 text-3xl font-bold">Frequently asked questions</h2>
      <p class="text-muted-foreground">Got questions? We have answers.</p>
    </div>

    <div class="mx-auto max-w-3xl space-y-6">
      {#each faqs as faq}
        <Card>
          <CardHeader>
            <CardTitle class="text-xl">{faq.question}</CardTitle>
          </CardHeader>
          <CardContent>
            <p class="text-muted-foreground">{faq.answer}</p>
          </CardContent>
        </Card>
      {/each}
    </div>
  </div>
</section>

<!-- CTA Section -->
<section class="bg-primary py-24 text-primary-foreground">
  <div class="container text-center">
    <h2 class="mb-4 text-3xl font-bold">Ready to transform your digital memory?</h2>
    <p class="mx-auto mb-8 max-w-2xl">
      Join thousands of users who have revolutionized how they organize their online world.
    </p>
    <div class="flex flex-wrap justify-center gap-4">
      <Button size="lg" variant="secondary">Get started for free</Button>
      <Button size="lg" variant="outline" class="bg-transparent hover:bg-primary-foreground/10">Book a demo</Button>
    </div>
  </div>
</section>

<!-- Footer -->
<footer class="border-t py-12">
  <div class="container">
    <div class="grid grid-cols-2 gap-8 md:grid-cols-4">
      <div>
        <h3 class="mb-4 text-lg font-semibold">Smaraka</h3>
        <p class="text-muted-foreground">Making digital memories accessible and meaningful.</p>
      </div>
      <div>
        <h3 class="mb-4 text-lg font-semibold">Product</h3>
        <ul class="space-y-2">
          <li><a href="#" class="text-muted-foreground hover:text-primary">Features</a></li>
          <li><a href="#" class="text-muted-foreground hover:text-primary">Pricing</a></li>
          <li><a href="#" class="text-muted-foreground hover:text-primary">Integrations</a></li>
          <li><a href="#" class="text-muted-foreground hover:text-primary">Roadmap</a></li>
        </ul>
      </div>
      <div>
        <h3 class="mb-4 text-lg font-semibold">Resources</h3>
        <ul class="space-y-2">
          <li><a href="#" class="text-muted-foreground hover:text-primary">Blog</a></li>
          <li><a href="#" class="text-muted-foreground hover:text-primary">Documentation</a></li>
          <li><a href="#" class="text-muted-foreground hover:text-primary">Help Center</a></li>
          <li><a href="#" class="text-muted-foreground hover:text-primary">Community</a></li>
        </ul>
      </div>
      <div>
        <h3 class="mb-4 text-lg font-semibold">Company</h3>
        <ul class="space-y-2">
          <li><a href="#" class="text-muted-foreground hover:text-primary">About</a></li>
          <li><a href="#" class="text-muted-foreground hover:text-primary">Careers</a></li>
          <li><a href="#" class="text-muted-foreground hover:text-primary">Contact</a></li>
          <li><a href="#" class="text-muted-foreground hover:text-primary">Legal</a></li>
        </ul>
      </div>
    </div>
    <div class="mt-12 border-t pt-8 text-center text-sm text-muted-foreground">
      <p>© {new Date().getFullYear()} Smaraka. All rights reserved.</p>
    </div>
  </div>
</footer>
