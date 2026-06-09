import { describe, it, expect } from "vitest";
import { render } from "@testing-library/react";
import EventCardSkeleton from "./EventCardSkeleton";

describe("EventCardSkeleton", () => {
  it("renders without crashing", () => {
    const { container } = render(<EventCardSkeleton />);
    expect(container.firstChild).toBeInTheDocument();
  });

  it("has the animate-pulse class for loading shimmer", () => {
    const { container } = render(<EventCardSkeleton />);
    expect(container.firstChild).toHaveClass("animate-pulse");
  });

  it("renders multiple placeholder blocks", () => {
    const { container } = render(<EventCardSkeleton />);
    // The card has the outer wrapper + several inner divs
    const divs = container.querySelectorAll("div");
    expect(divs.length).toBeGreaterThan(3);
  });
});
