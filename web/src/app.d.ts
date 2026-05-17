declare global {
  namespace App {
    interface Locals {
      user?: { id: number }
    }
    interface PageData {
      user?: { id: number }
    }
  }
}

export {};
