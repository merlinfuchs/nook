import HomeCTABanner from "@/components/home/HomeCTABanner";
import HomeFAQ from "@/components/home/HomeFAQ";
import HomeFeatures from "@/components/home/HomeFeatures";
import HomeFooter from "@/components/home/HomeFooter";
import HomeHero from "@/components/home/HomeHero";
import HomeModules from "@/components/home/HomeModules";
import HomeNavBar from "@/components/home/HomeNavBar";

export default function Home() {
  return (
    <>
      <HomeNavBar />
      <main className="pt-16 xs:pt-20 sm:pt-24">
        <HomeHero />
        <HomeFeatures />
        <HomeFAQ />
        <HomeModules />
        <HomeCTABanner />
        <HomeFooter />
      </main>
    </>
  );
}
