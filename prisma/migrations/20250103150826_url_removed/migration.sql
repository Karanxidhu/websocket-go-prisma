/*
  Warnings:

  - You are about to drop the column `url` on the `MediaFile` table. All the data in the column will be lost.

*/
-- DropIndex
DROP INDEX "MediaFile_url_key";

-- AlterTable
ALTER TABLE "MediaFile" DROP COLUMN "url";
