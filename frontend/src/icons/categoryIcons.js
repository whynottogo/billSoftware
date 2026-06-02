var CATEGORY_ICON_MAP = {
  "food-bowl":           { color: "#FF6B6B", icon: "Bowl" },
  "food-utensils":       { color: "#FF8E53", icon: "ForkSpoon" },
  "shopping-bag":        { color: "#FFA940", icon: "ShoppingBag" },
  "daily-home":          { color: "#FFD93D", icon: "HomeFilled" },
  "drink-cup":           { color: "#C084FC", icon: "MilkTea" },
  "transport-car":       { color: "#60A5FA", icon: "Van" },
  "food-carrot":         { color: "#34D399", icon: "Food" },
  "food-apple":          { color: "#F472B6", icon: "Apple" },
  "food-candy":          { color: "#FB923C", icon: "Lollipop" },
  "sport-run":           { color: "#4ADE80", icon: "Basketball" },
  "entertainment-game":  { color: "#A78BFA", icon: "Film" },
  "comm-phone":          { color: "#38BDF8", icon: "Iphone" },
  "clothes-tshirt":      { color: "#F472B6", icon: "ShoppingCartFull" },
  "beauty-mirror":       { color: "#EC4899", icon: "Brush" },
  "housing-key":         { color: "#FBBF24", icon: "Key" },
  "home-sofa":           { color: "#FB923C", icon: "House" },
  "family-child":        { color: "#F472B6", icon: "IceCream" },
  "family-elder":        { color: "#EF4444", icon: "Avatar" },
  "social-handshake":    { color: "#60A5FA", icon: "ChatDotRound" },
  "travel-plane":        { color: "#38BDF8", icon: "Location" },
  "drink-wine":          { color: "#A78BFA", icon: "GobletFull" },
  "digital-phone":       { color: "#64748B", icon: "Monitor" },
  "transport-wheel":     { color: "#64748B", icon: "Bicycle" },
  "medical-cross":       { color: "#EF4444", icon: "FirstAidKit" },
  "edu-book":            { color: "#3B82F6", icon: "Reading" },
  "pet-paw":             { color: "#FB923C", icon: "Cherry" },
  "gift-box":            { color: "#EC4899", icon: "Present" },
  "work-briefcase":      { color: "#64748B", icon: "Briefcase" },
  "repair-wrench":       { color: "#64748B", icon: "Tools" },
  "donate-heart":        { color: "#EF4444", icon: "Coin" },
  "lottery-star":        { color: "#FBBF24", icon: "StarFilled" },
  "family-people":       { color: "#60A5FA", icon: "User" },
  "delivery-package":    { color: "#FB923C", icon: "Box" },
  "income-salary":       { color: "#22C55E", icon: "Money" },
  "income-parttime":     { color: "#3B82F6", icon: "Opportunity" },
  "income-invest":       { color: "#F59E0B", icon: "TrendCharts" },
  "income-gift":         { color: "#EF4444", icon: "Present" },
  "income-other":        { color: "#94A3B8", icon: "MoreFilled" }
};

var CATEGORY_ICON_KEYS = Object.keys(CATEGORY_ICON_MAP);

function getCategoryIcon(iconKey) {
  if (iconKey && CATEGORY_ICON_MAP[iconKey]) {
    return CATEGORY_ICON_MAP[iconKey];
  }
  return null;
}

function getAllCategoryIcons() {
  return CATEGORY_ICON_KEYS.map(function(key) {
    return {
      key: key,
      color: CATEGORY_ICON_MAP[key].color,
      icon: CATEGORY_ICON_MAP[key].icon
    };
  });
}

export { CATEGORY_ICON_MAP, CATEGORY_ICON_KEYS, getCategoryIcon, getAllCategoryIcons };
