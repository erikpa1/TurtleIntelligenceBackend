package llmCtrl

import (
	"encoding/json"
	"fmt"
	"github.com/erikpa1/turtle/llm/llmModels"
	"regexp"
	"sort"
	"strings"
)

// ContentSegment represents a segment with its position and type
type ContentSegment struct {
	Start int
	End   int
	Type  string
	Text  string
}

// isJSON checks if text is valid JSON
func isJSON(text string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(text), &js) == nil
}

// isHTML checks if text contains HTML tags
func isHTML(text string) bool {
	htmlPattern := `<[^>]+>`
	matched, _ := regexp.MatchString(htmlPattern, text)
	return matched
}

// isValidHTML checks if HTML content is complete and renderable
func isValidHTML(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}

	// Must contain HTML tags
	if !isHTML(text) {
		return false
	}

	// Should be a complete HTML structure or valid HTML fragment
	// Check for basic HTML document structure or valid HTML fragment
	hasHtmlTag := regexp.MustCompile(`(?i)<html[^>]*>`).MatchString(text)
	hasBodyTag := regexp.MustCompile(`(?i)<body[^>]*>`).MatchString(text)
	hasValidTags := regexp.MustCompile(`(?i)<(div|section|article|main|header|footer|nav|aside|p|h[1-6]|ul|ol|li|table|form|input|button|span|a|img)[^>]*>`).MatchString(text)

	return hasHtmlTag || hasBodyTag || hasValidTags
}

// extractHTMLFromCodeBlocks extracts HTML content from markdown code blocks
func extractHTMLFromCodeBlocks(text string) []string {
	var segments []string

	// Improved patterns to match ```html ... ``` blocks
	patterns := []string{
		// Match ```html with content and closing ```
		`(?s)` + "`" + `{3}html\s*\n(.*?)` + "`" + `{3}`,
		// Match ```HTML (case insensitive)
		`(?s)` + "`" + `{3}HTML\s*\n(.*?)` + "`" + `{3}`,
		// Match ``` followed by HTML-like content (fallback)
		`(?s)` + "`" + `{3}\s*\n(<(?:!DOCTYPE\s+html|html|body|div|section|article|main|header|footer)[^>]*>.*?)` + "`" + `{3}`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(text, -1)

		for _, match := range matches {
			if len(match) > 1 {
				content := strings.TrimSpace(match[1])
				if content != "" {
					segments = append(segments, content)
				}
			}
		}
	}

	return segments
}

// extractFromCodeBlocks extracts content from markdown code blocks (for JSON)
func extractFromCodeBlocks(text string, language string) []string {
	var segments []string

	// Pattern for markdown code blocks with specific language
	pattern := fmt.Sprintf(`(?s)`+"`"+`{3}%s\s*\n(.*?)`+"`"+`{3}`, language)
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) > 1 {
			content := strings.TrimSpace(match[1])
			if content != "" {
				segments = append(segments, content)
			}
		}
	}

	return segments
}

// repairJSON attempts to fix common JSON issues
func repairJSON(text string) string {
	// Remove leading/trailing whitespace
	text = strings.TrimSpace(text)

	// Fix common issues
	text = strings.ReplaceAll(text, "'", "\"") // Replace single quotes with double quotes
	text = strings.ReplaceAll(text, "True", "true")
	text = strings.ReplaceAll(text, "False", "false")
	text = strings.ReplaceAll(text, "None", "null")
	text = strings.ReplaceAll(text, "Null", "null")

	// Remove trailing commas before closing brackets/braces
	re := regexp.MustCompile(`,(\s*[}\]])`)
	text = re.ReplaceAllString(text, "$1")

	// Try to fix missing quotes around keys
	keyPattern := regexp.MustCompile(`(\w+)(\s*:\s*)`)
	text = keyPattern.ReplaceAllString(text, `"$1"$2`)

	// Fix duplicate quotes
	text = regexp.MustCompile(`""+`).ReplaceAllString(text, `"`)

	return text
}

// repairHTML attempts to fix common HTML issues and ensure it's frontend-ready
func repairHTML(text string) string {
	// Remove leading/trailing whitespace
	text = strings.TrimSpace(text)

	// Fix common HTML entity issues
	text = strings.ReplaceAll(text, "&amp;amp;", "&amp;")

	// If we have an opening body tag but no closing tag, try to add it
	if regexp.MustCompile(`(?i)<body[^>]*>`).MatchString(text) && !regexp.MustCompile(`(?i)</body>`).MatchString(text) {
		text = text + "\n</body>"
	}

	return text
}

// findContentSegments finds all structured content segments with their positions
func findContentSegments(text string) []ContentSegment {
	var segments []ContentSegment

	// Find JSON segments with positions
	segments = append(segments, findJSONSegmentsWithPositions(text)...)

	// Find HTML segments with positions
	segments = append(segments, findHTMLSegmentsWithPositions(text)...)

	// Sort segments by start position
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].Start < segments[j].Start
	})

	// Remove overlapping segments (keep the first one found)
	var cleanSegments []ContentSegment
	for i, segment := range segments {
		overlap := false
		for j := 0; j < i; j++ {
			if segments[j].Start <= segment.Start && segment.Start < segments[j].End {
				overlap = true
				break
			}
			if segment.Start <= segments[j].Start && segments[j].Start < segment.End {
				overlap = true
				break
			}
		}
		if !overlap {
			cleanSegments = append(cleanSegments, segment)
		}
	}

	return cleanSegments
}

// findJSONSegmentsWithPositions finds JSON segments and their positions
func findJSONSegmentsWithPositions(text string) []ContentSegment {
	var segments []ContentSegment

	// First, try to extract from code blocks
	codeBlockPattern := `(?s)` + "`" + `{3}json\s*\n(.*?)` + "`" + `{3}`
	re := regexp.MustCompile(codeBlockPattern)
	matches := re.FindAllStringSubmatchIndex(text, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			fullStart, fullEnd := match[0], match[1]
			contentStart, contentEnd := match[2], match[3]
			content := text[contentStart:contentEnd]
			repaired := repairJSON(content)

			if isJSON(repaired) {
				segments = append(segments, ContentSegment{
					Start: fullStart,
					End:   fullEnd,
					Type:  "json",
					Text:  repaired,
				})
			}
		}
	}

	// If found in code blocks, prioritize those
	if len(segments) > 0 {
		return segments
	}

	// Look for JSON patterns in the text
	jsonPatterns := []string{
		`\{(?:[^{}]|\{[^{}]*\})*\}`,     // Nested objects (limited depth)
		`\[(?:[^\[\]]|\[[^\[\]]*\])*\]`, // Nested arrays (limited depth)
		`\{[^{}]*\}`,                    // Simple objects
		`\[[^\[\]]*\]`,                  // Simple arrays
	}

	for _, pattern := range jsonPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringIndex(text, -1)

		for _, match := range matches {
			if len(match) >= 2 {
				content := text[match[0]:match[1]]
				repaired := repairJSON(content)

				if isJSON(repaired) {
					segments = append(segments, ContentSegment{
						Start: match[0],
						End:   match[1],
						Type:  "json",
						Text:  repaired,
					})
				}
			}
		}
	}

	return segments
}

// findHTMLSegmentsWithPositions finds HTML segments and their positions
func findHTMLSegmentsWithPositions(text string) []ContentSegment {
	var segments []ContentSegment

	// First, try to extract from markdown code blocks
	htmlCodeBlockPatterns := []string{
		`(?s)` + "`" + `{3}html\s*\n(.*?)` + "`" + `{3}`,
		`(?s)` + "`" + `{3}HTML\s*\n(.*?)` + "`" + `{3}`,
		`(?s)` + "`" + `{3}\s*\n(<(?:!DOCTYPE\s+html|html|body|div|section|article|main|header|footer)[^>]*>.*?)` + "`" + `{3}`,
	}

	for _, pattern := range htmlCodeBlockPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatchIndex(text, -1)

		for _, match := range matches {
			if len(match) >= 4 {
				fullStart, fullEnd := match[0], match[1]
				contentStart, contentEnd := match[2], match[3]
				content := text[contentStart:contentEnd]
				repaired := repairHTML(content)

				if isValidHTML(repaired) {
					segments = append(segments, ContentSegment{
						Start: fullStart,
						End:   fullEnd,
						Type:  "html",
						Text:  repaired,
					})
				}
			}
		}
	}

	// If found valid HTML in code blocks, prioritize those
	if len(segments) > 0 {
		return segments
	}

	// Look for HTML patterns in the text
	htmlPatterns := []string{
		`(?i)<body[^>]*>.*?</body>`,
		`(?i)<html[^>]*>.*?</html>`,
		`(?i)<head[^>]*>.*?</head>`,
		`(?i)<div[^>]*>.*?</div>`,
		`(?i)<section[^>]*>.*?</section>`,
		`(?i)<article[^>]*>.*?</article>`,
		`(?i)<main[^>]*>.*?</main>`,
		`(?i)<header[^>]*>.*?</header>`,
		`(?i)<footer[^>]*>.*?</footer>`,
		`(?i)<nav[^>]*>.*?</nav>`,
		`(?i)<aside[^>]*>.*?</aside>`,
	}

	for _, pattern := range htmlPatterns {
		re := regexp.MustCompile("(?s)" + pattern)
		matches := re.FindAllStringIndex(text, -1)

		for _, match := range matches {
			if len(match) >= 2 {
				content := text[match[0]:match[1]]
				repaired := repairHTML(content)

				if isValidHTML(repaired) {
					segments = append(segments, ContentSegment{
						Start: match[0],
						End:   match[1],
						Type:  "html",
						Text:  repaired,
					})
				}
			}
		}
	}

	// Also look for incomplete body tags
	// Look for body opening tags and check if they have matching closing tags
	bodyOpenPattern := `(?i)<body[^>]*>`
	bodyOpenRe := regexp.MustCompile(bodyOpenPattern)
	bodyOpenMatches := bodyOpenRe.FindAllStringIndex(text, -1)

	for _, match := range bodyOpenMatches {
		if len(match) >= 2 {
			// Extract text from body tag to end of string
			remainingText := text[match[0]:]

			// Check if there's a closing </body> tag
			if !regexp.MustCompile(`(?i)</body>`).MatchString(remainingText) {
				// This is an incomplete body tag - take content to end of text
				content := remainingText
				repaired := repairHTML(content)

				if isValidHTML(repaired) {
					segments = append(segments, ContentSegment{
						Start: match[0],
						End:   len(text),
						Type:  "html",
						Text:  repaired,
					})
				}
			}
		}
	}

	return segments
}

// Legacy functions for backward compatibility
func findJSONSegments(text string) []string {
	segments := findJSONSegmentsWithPositions(text)
	var results []string
	for _, segment := range segments {
		results = append(results, segment.Text)
	}
	return results
}

func findHTMLSegments(text string) []string {
	segments := findHTMLSegmentsWithPositions(text)
	var results []string
	for _, segment := range segments {
		results = append(results, segment.Text)
	}
	return results
}

// FindJSON detects all JSON segments in text
func FindJSON(text string) []llmModels.ContentBlock {
	var results []llmModels.ContentBlock

	segments := findJSONSegments(text)
	for _, segment := range segments {
		results = append(results, llmModels.ContentBlock{Type: "json", Text: segment})
	}

	return results
}

// FindHTML detects all HTML segments in text, prioritizing code blocks
func FindHTML(text string) []llmModels.ContentBlock {
	var results []llmModels.ContentBlock

	segments := findHTMLSegments(text)
	for _, segment := range segments {
		results = append(results, llmModels.ContentBlock{Type: "html", Text: segment})
	}

	return results
}

// FindBodyHTML specifically detects content with body tags
func FindBodyHTML(text string) []llmModels.ContentBlock {
	var results []llmModels.ContentBlock

	// First check code blocks for body content
	codeBlockSegments := extractHTMLFromCodeBlocks(text)
	for _, segment := range codeBlockSegments {
		if regexp.MustCompile(`(?i)<body[^>]*>`).MatchString(segment) {
			repaired := repairHTML(segment)
			if isValidHTML(repaired) {
				results = append(results, llmModels.ContentBlock{Type: "html", Text: repaired})
			}
		}
	}

	if len(results) > 0 {
		return results
	}

	bodyPattern := `(?i)<body[^>]*>.*?</body>`
	bodyRe := regexp.MustCompile("(?s)" + bodyPattern)
	bodyMatches := bodyRe.FindAllString(text, -1)

	for _, match := range bodyMatches {
		repaired := repairHTML(match)
		if isValidHTML(repaired) {
			results = append(results, llmModels.ContentBlock{Type: "html", Text: repaired})
		}
	}

	// Also look for incomplete body tags
	incompleteBodyPattern := `(?i)<body[^>]*>(?:(?!</body>).)*$`
	incompleteRe := regexp.MustCompile("(?s)" + incompleteBodyPattern)
	incompleteMatches := incompleteRe.FindAllString(text, -1)

	for _, match := range incompleteMatches {
		repaired := repairHTML(match)
		if isValidHTML(repaired) {
			results = append(results, llmModels.ContentBlock{Type: "html", Text: repaired})
		}
	}

	return results
}

// FindAllContent detects all types of content segments in text and splits into alternating segments
func FindAllContent(text string) []llmModels.ContentBlock {
	var results []llmModels.ContentBlock

	// Handle empty text
	if strings.TrimSpace(text) == "" {
		return results
	}

	// Find all structured content segments with positions
	structuredSegments := findContentSegments(text)

	// If no structured content found, return entire text as one block
	if len(structuredSegments) == 0 {
		results = append(results, llmModels.ContentBlock{Type: "text", Text: text})
		return results
	}

	// Build alternating segments of text and structured content
	currentPos := 0

	for _, segment := range structuredSegments {
		// Add text segment before structured content (if any)
		if segment.Start > currentPos {
			textContent := text[currentPos:segment.Start]
			if strings.TrimSpace(textContent) != "" {
				results = append(results, llmModels.ContentBlock{
					Type: "text",
					Text: textContent,
				})
			}
		}

		// Add the structured content segment
		results = append(results, llmModels.ContentBlock{
			Type: segment.Type,
			Text: segment.Text,
		})

		currentPos = segment.End
	}

	// Add remaining text after last structured segment (if any)
	if currentPos < len(text) {
		textContent := text[currentPos:]
		if strings.TrimSpace(textContent) != "" {
			results = append(results, llmModels.ContentBlock{
				Type: "text",
				Text: textContent,
			})
		}
	}

	return results
}
